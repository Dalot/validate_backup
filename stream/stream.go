package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

type List map[string][]ParsedObj

type Stream struct {
	beforeListMsg chan *Obj
	afterListMsg  chan *Obj
	compMsgs      chan Comparison
	list          List
	mutex         *sync.RWMutex
}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	var wgParse sync.WaitGroup
	var wgCompare sync.WaitGroup
	var wgResult sync.WaitGroup
	s := &Stream{}
	s.mutex = &sync.RWMutex{}
	s.beforeListMsg = make(chan *Obj)
	s.afterListMsg = make(chan *Obj)
	s.compMsgs = make(chan Comparison)
	s.list = map[string][]ParsedObj{}

	wgParse.Add(2)
	go s.Parse(beforeReader, s.beforeListMsg, &wgParse)
	go s.Parse(afterReader, s.afterListMsg, &wgParse)
	go func() {
		wgParse.Wait()
	}()

	wgCompare.Add(2)
	go s.comparisonWorker(&wgCompare, s.beforeListMsg, Before)
	go s.comparisonWorker(&wgCompare, s.afterListMsg, After)

	wgResult.Add(1)
	go func(wgResult *sync.WaitGroup) {
		wgCompare.Wait()

		for _, list := range s.list {
			len := len(list)
			if len != 2 {

				s.compMsgs <- Different

				return
			} else {
				log.Fatalln(list)
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
			if obj.Id == "" {
				log.Println("empty id: ", err)
				s.compMsgs <- Different
				return
			}
			log.Fatal("[2]", err)
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
		s.list[msg.Id] = append(s.list[msg.Id], parsedObj)
		s.mutex.Unlock()

		s.mutex.RLock()
		len := len(s.list[msg.Id])
		hasTwoItems := len == 2

		if hasTwoItems {
			differentTypes := s.list[msg.Id][0].Type == s.list[msg.Id][1].Type
			if differentTypes {
				s.compMsgs <- DuplicateIds
				s.mutex.RUnlock()
				return
			}

			namesAreDifferent := s.list[msg.Id][0].Name != s.list[msg.Id][1].Name
			if namesAreDifferent {
				s.compMsgs <- Different
			}
			s.mutex.RUnlock()
			s.mutex.Lock()
			delete(s.list, msg.Id) // Whether Objects are equal ot Different, we can free memory
			s.mutex.Unlock()

		} else {
			s.mutex.RUnlock()
		}

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
