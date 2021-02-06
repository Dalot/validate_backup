package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
)

type Comparison int

const (
	Equal Comparison = iota
	InvalidJson
	DuplicateIds
	Different
	None
)

type Obj struct {
	Id   string
	Name string
}

type List []*Obj

var lists = make(map [int]*List)

type Stream struct{}

type Tracker struct {
	parentIndex int
	childIndex int
}


// afterBackupIDMap is a map that will point the id to an index for the AfterBackup List
var afterBackupIDMap = map[string]Tracker{}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	stream := &Stream{}

	listOfListsBefore, comp := stream.Parse(beforeReader, nil)
	if comp != None {
		return stream, comp
	}
	listOfListsAfter, comp := stream.Parse(afterReader, afterBackupIDMap)
	if comp != None {
		return stream, comp
	}

	//result := stream.compare(listBeforeBackup, listAfterBackup)
	result := stream.compare(listOfListsBefore, listOfListsAfter)

	return stream, result
}

func (s *Stream) compare(listBeforeBackup map [int]*List, listAfterBackup map [int]*List) Comparison {

	equalFoundObjs := 0
	for i := range listBeforeBackup {
		for _, obj := range *listBeforeBackup[i] {
			if tracker, exists := afterBackupIDMap[obj.Id]; !exists {
				log.Println("Did not find id(", obj.Id, ") in the file after the backup")
				return Different
			} else {
				listOfAfter := *listAfterBackup[tracker.parentIndex]
				otherObj := listOfAfter[tracker.childIndex]
				if obj.Name == otherObj.Name {
					//log.Print("Found id(", obj.Data.Id, ") in both files with the same data")
					//log.Print("Location: beforeBackup index:", beforeBackupIDMap[obj.Data.Id], ", afterBackup:", index)
					equalFoundObjs++
				} else {
					return Different
				}
			}
		}
	}

	//log.Println("equalFoundObjs: (", equalFoundObjs, ")")

	return Equal
}

// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) Parse(input io.Reader, idMap map[string]Tracker) (map [int]*List, Comparison) {
	var count int = 0
	var step = 100000
	var stepIndex = 1
	decoder := json.NewDecoder(input)

	// read open bracket
	t, err := decoder.Token()
	if err != nil {
		log.Fatal("[1]", err)
	}
	fmt.Printf("%T: %v\n", t, t)
	
	// while the array contains values
	for decoder.More() {
		var data Obj
		// decode an array value (Message)
		err := decoder.Decode(&data)

		if err != nil {
			log.Fatal("[2]", err)
		}

		obj := Obj{
			Name:  data.Name,
			Id:  data.Id,
		}

		if idMap != nil {
			if tracker, exists := idMap[data.Id]; exists {
				log.Println("IDs should be unique. The id:", data.Id, " is already present")
				log.Println("The duplicate ID has parentIndex of:", tracker.parentIndex, " and a childIndex of:", tracker.childIndex)
				return make(map [int]*List,0), DuplicateIds
			}
		}

		currentStepIndex := math.Ceil(float64(count + 1) / float64(step) ) 
		if  currentStepIndex > float64(stepIndex)  {
			stepIndex++
		}
		currentCount := count - step * (stepIndex - 1)

		if _, exists := lists[int(currentStepIndex)]; !exists {
			lists[int(currentStepIndex)] = &List{}
		}
		
		*lists[int(currentStepIndex)]= append(*lists[int(currentStepIndex)], &obj)
		if idMap != nil {
			t := Tracker{
				parentIndex: int(currentStepIndex),
				childIndex: currentCount,
			}
			idMap[data.Id] = t
		}
		count++
		
		//if count == 7000000 {
		//	fmt.Println("Real size of list:", unsafe.Sizeof(list)+unsafe.Sizeof([7000000]*Obj{}))
		//}
	}

	// read closing bracket
	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
	fmt.Println("--------")
	fmt.Println("count : = ", count)
	fmt.Println("--------")

	return lists, None
}

func (s *Stream) Flush() {
	afterBackupIDMap = map[string]Tracker{}
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
