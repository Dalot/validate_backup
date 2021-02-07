package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"sync"
)

type Comparison int

const (
	Equal Comparison = iota
	InvalidJson
	DuplicateIds
	Different
	None
)

type Type string

const (
	Before Type = "Before"
	After       = "After"
)

type Data struct {
	Id   string
	Name string
}

type Obj struct {
	Id   string
	Name string
}

type ParsedObj struct {
	Id   string
	Name string
	Type Type
}

type List map[string][]*ParsedObj

type JsonObj *map[string]interface{}

type ListInterface map[string][]*map[string]interface{}

type Stream struct {
	beforeListMsg          chan *Obj
	afterListMsg           chan *Obj
	beforeListMsgInterface chan *map[string]interface{}
	afterListMsgInterface  chan *map[string]interface{}
	compMsgs               chan Comparison
	list                   *List
	listInterface          *ListInterface
	mutex                  *sync.RWMutex
}

func New() *Stream {
	s := &Stream{}
	s.mutex = &sync.RWMutex{}
	s.beforeListMsg = make(chan *Obj)
	s.afterListMsg = make(chan *Obj)
	s.beforeListMsgInterface = make(chan *map[string]interface{})
	s.afterListMsgInterface = make(chan *map[string]interface{})
	s.compMsgs = make(chan Comparison)
	s.list = &List{}
	s.listInterface = &ListInterface{}

	return s
}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	var wgParse sync.WaitGroup
	var wgCompare sync.WaitGroup
	var wgResult sync.WaitGroup
	s := New()

	wgParse.Add(2)
	//go s.Parse(beforeReader, s.beforeListMsg, &wgParse)
	//go s.Parse(afterReader, s.afterListMsg, &wgParse)
	go s.ParseInterface(beforeReader, s.beforeListMsgInterface, &wgParse)
	go s.ParseInterface(afterReader, s.afterListMsgInterface, &wgParse)
	go func() {
		wgParse.Wait()
	}()

	wgCompare.Add(2)
	//go s.comparisonWorker(&wgCompare, s.beforeListMsg, Before)
	//go s.comparisonWorker(&wgCompare, s.afterListMsg, After)
	go s.comparisonWorkerInterface(&wgCompare, s.beforeListMsgInterface, Before)
	go s.comparisonWorkerInterface(&wgCompare, s.afterListMsgInterface, After)

	wgResult.Add(1)
	go func(wgResult *sync.WaitGroup) {
		wgCompare.Wait()

		for _, list := range *s.listInterface {
			len := len(list)
			if len != 2 {

				s.compMsgs <- Different

				return
			} else {
				log.Fatalln("Something happened")
			}
		}
		close(s.compMsgs)
		wgResult.Done()
	}(&wgResult)

	go func() {
		wgResult.Wait()
	}()

	for comp := range s.compMsgs {
		if comp != Equal {
			log.Println("comp:", comp)
			return s, comp
		} 
	}

	log.Println("FINISHED")
	return s, Equal
}

// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) Parse(input io.Reader, msgs chan<- *Obj, wg *sync.WaitGroup) {
	defer wg.Done()
	var count uint64 = 0
	decoder := json.NewDecoder(input)

	// read open bracket
	t, err := decoder.Token()
	if err != nil {
		log.Fatal("[1]", err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for decoder.More() {
		var obj Obj
		// decode an array value (Message)
		err := decoder.Decode(&obj)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && obj.Id == "" {
				log.Println("empty id: ", err)
				s.compMsgs <- InvalidJson
				return
			} else {
				log.Fatal("[2]", err)
			}
		}

		if obj.Id == "" {
			log.Println("empty id: ", err)
			s.compMsgs <- InvalidJson
			return
		}

		msgs <- &obj

		count++
	}

	// read closing bracket
	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
	fmt.Println("--------")
	fmt.Println("count : = " + strconv.FormatUint(count, 10))
	fmt.Println("--------")
	defer close(msgs)
}

// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) ParseInterface(input io.Reader, msgs chan<- *map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var count uint64 = 0
	decoder := json.NewDecoder(input)

	// read open bracket
	t, err := decoder.Token()
	if err != nil {
		log.Fatal("[1]", err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for decoder.More() {
		//var obj Obj
		var data map[string]interface{}
		// decode an array value (Message)
		err := decoder.Decode(&data)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && (data["id"] == "" || data == nil) {
				log.Println("empty id: ", err)
				s.compMsgs <- InvalidJson
				return
			}

			log.Fatal("[2]", err)

		}

		if len(data) == 0 {
			log.Println("empty id: ", err)
			s.compMsgs <- InvalidJson
			return
		}

		if data["id"] == "" {
			log.Println("empty id: ", err)
			s.compMsgs <- InvalidJson
			return
		}

		msgs <- &data

		count++
	}

	// read closing bracket
	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
	fmt.Println("--------")
	fmt.Println("count : = " + strconv.FormatUint(count, 10))
	fmt.Println("--------")
	defer close(msgs)
}

func (s *Stream) comparisonWorker(wg *sync.WaitGroup, objMsgs chan *Obj, msgType Type) {

	defer wg.Done()

	for msg := range objMsgs {
		parsedObj := ParsedObj{
			Id:   msg.Id,
			Name: msg.Name,
			Type: msgType,
		}
		if len(parsedObj.Id) == 0 {
			s.compMsgs <- Different
			return
		}

		s.mutex.Lock()
		(*s.list)[msg.Id] = append((*s.list)[msg.Id], &parsedObj)
		s.mutex.Unlock()

		s.mutex.RLock()
		len := len((*s.list)[msg.Id])
		hasTwoItems := len == 2

		if hasTwoItems {
			differentTypes := (*s.list)[msg.Id][0].Type == (*s.list)[msg.Id][1].Type
			if differentTypes {
				s.compMsgs <- DuplicateIds
				s.mutex.RUnlock()
				return
			}

			namesAreDifferent := (*s.list)[msg.Id][0].Name != (*s.list)[msg.Id][1].Name
			if namesAreDifferent {
				s.compMsgs <- Different
			}
			s.mutex.RUnlock()
			s.mutex.Lock()
			delete(*s.list, msg.Id) // Whether Objects are equal ot Different, we can free memory
			s.mutex.Unlock()

		} else {
			s.mutex.RUnlock()
		}

	}

}

func (s *Stream) comparisonWorkerInterface(wg *sync.WaitGroup, objMsgs chan *map[string]interface{}, msgType Type) {

	defer wg.Done()

	for msg := range objMsgs {

		id, ok := (*msg)["id"].(string)
		if !ok {
			log.Fatalln(ok)
		}
		if len(id) == 0 {
			s.compMsgs <- Different
			return
		}

		s.mutex.Lock()
		(*msg)["type"] = msgType
		(*s.listInterface)[id] = append((*s.listInterface)[id], msg)
		

		
		length := len((*s.listInterface)[id])
		hasTwoItems := length == 2

		if hasTwoItems {
			sameList := (*(*s.listInterface)[id][0])["type"] == (*(*s.listInterface)[id][1])["type"]
			if sameList {
				s.mutex.Unlock()
				s.compMsgs <- DuplicateIds
				return
			}
			
			
			newObjects := []map[string]string{}
			for _, obj := range (*s.listInterface)[id] {
				delete(*obj, "type")
				newObj := make(map[string]string)
				for key, value := range *obj {
					var newValue string
					switch v := value.(type) {
					case nil:
						newValue = ""
					case int:
						newValue = strconv.Itoa(value.(int))
					case bool:
						newValue = strconv.FormatBool(value.(bool))
					case string:
						newValue = value.(string)
					//case Type:
					//	newValue = string(value.(Type))
					default:
						log.Fatalln("could not convert value: ", v)
					}

					newObj[key] = newValue
				}
				newObjects = append(newObjects, newObj)
			}
			
			eq := reflect.DeepEqual(newObjects[0], newObjects[1])
			if !eq {
				s.compMsgs <- Different
				
				delete(*s.listInterface, id) // Whether Objects are equal ot Different, we can free memory
				s.mutex.Unlock()
				return
			} 

			
			delete(*s.listInterface, id) // Whether Objects are equal ot Different, we can free memory
		} 
		s.mutex.Unlock()

	}

}

func (c Comparison) String() string {
	switch c {
	case Equal:
		return "Files have equal objects"
	case InvalidJson:
		return "The json seems invalid"
	case DuplicateIds:
		return "There seems to be duplicate IDs"
	case Different:
		return "Files have different objects"
	default:
		return fmt.Sprintf("%d", int(c))
	}
}
