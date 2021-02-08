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

type JsonObj *map[string]interface{}

type List map[string][]*map[string]interface{}

type Stream struct {
	beforeListMsgs chan *map[string]interface{}
	afterListMsgs  chan *map[string]interface{}
	compMsgs               chan Comparison
	list          *List
	mutex                  *sync.RWMutex
}

func New() *Stream {
	s := &Stream{}
	s.mutex = &sync.RWMutex{}
	s.beforeListMsgs = make(chan *map[string]interface{})
	s.afterListMsgs = make(chan *map[string]interface{})
	s.compMsgs = make(chan Comparison)
	s.list = &List{}

	return s
}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	var wgParse sync.WaitGroup
	var wgCompare sync.WaitGroup
	var wgResult sync.WaitGroup
	s := New()

	wgParse.Add(2)
	go s.ParseInterface(beforeReader, s.beforeListMsgs, &wgParse)
	go s.ParseInterface(afterReader, s.afterListMsgs, &wgParse)
	go func() {
		wgParse.Wait()
	}()

	// Tested with 500mb json file
	// 10 goroutines to each list messages is able to keep memory below 250mb
	// 30 goroutines seems to do a bit worse, going up to 400mb
	// 3 - 5 goroutines each seems to be the sweet deal, keeping below 100mb
	for i := 0; i < 5; i++ {
		wgCompare.Add(2)
		go s.comparisonWorkerInterface(&wgCompare, s.beforeListMsgs, Before)
		go s.comparisonWorkerInterface(&wgCompare, s.afterListMsgs, After)
	}

	wgResult.Add(1)
	go func(wgResult *sync.WaitGroup) {
		wgCompare.Wait()

		for _, list := range *s.list {
			len := len(list)
			if len != 2 {

				s.compMsgs <- Different

				return
			} else {
				log.Println("Something happened")
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

		val, ok := data["id"].(float64)
		if ok {
			data["id"] = strconv.Itoa(int(val))
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
		(*s.list)[id] = append((*s.list)[id], msg)
		
		length := len((*s.list)[id])
		hasTwoItems := length == 2

		if hasTwoItems {
			sameList := (*(*s.list)[id][0])["type"] == (*(*s.list)[id][1])["type"]
			if sameList {
				s.mutex.Unlock()
				s.compMsgs <- DuplicateIds
				return
			}
			
			
			newObjects := []map[string]string{}
			for _, obj := range (*s.list)[id] {
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
				
				delete(*s.list, id) // Whether Objects are equal ot Different, we can free memory
				s.mutex.Unlock()
				return
			} 

			
			delete(*s.list, id) // Whether Objects are equal ot Different, we can free memory
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
