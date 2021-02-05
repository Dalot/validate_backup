package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Comparison int

const (
	Equal Comparison = iota
	InvalidJson
	DuplicateIds
	Different
	None
)

type Data struct {
	Id   string
	Name string
}

type Obj struct {
	Data  Data
	bytes int64
	start int64 //byte offset start
	end   int64 // /byte offset end
}

type List []Obj

type Stream struct{}

// beforeBackupIDMap is a map that will point the id to an index for the BeforeBackup List
var beforeBackupIDMap = map[string]uint64{}

// afterBackupIDMap is a map that will point the id to an index for the AfterBackup List
var afterBackupIDMap = map[string]uint64{}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	stream := &Stream{}

	listBeforeBackup, comp := stream.Parse(beforeReader, beforeBackupIDMap)
	if comp != None {
		return stream, comp
	}
	listAfterBackup, comp := stream.Parse(afterReader, afterBackupIDMap)
	if comp != None {
		return stream, comp
	}

	result := stream.compare(listBeforeBackup, listAfterBackup)

	return stream, result
}

func (s *Stream) compare(listBeforeBackup *List, listAfterBackup *List) Comparison {

	equalFoundObjs := 0
	for _, obj := range *listBeforeBackup {
		if index, exists := afterBackupIDMap[obj.Data.Id]; !exists {
			log.Println("Did not find id(", obj.Data.Id, ") in the file after the backup")
			return Different
		} else {

			if obj.Data.Name == (*listAfterBackup)[index].Data.Name {
				log.Print("Found id(", obj.Data.Id, ") in both files with the same data")
				log.Print("Location: beforeBackup index:", beforeBackupIDMap[obj.Data.Id], ", afterBackup:", index)
				equalFoundObjs++
			} else {
				return Different
			}
		}
	}

	log.Println("equalFoundObjs: (", equalFoundObjs, ")")

	return Equal
}

// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) Parse(input io.Reader, idMap map[string]uint64) (*List, Comparison)  {
	var count uint64 = 0
	var list *List = &List{}
	decoder := json.NewDecoder(input)

	// read open bracket
	t, err := decoder.Token()
	if err != nil {
		log.Fatal("[1]", err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for decoder.More() {
		var data Data
		// decode an array value (Message)
		currentOffset := decoder.InputOffset()
		err := decoder.Decode(&data)

		if err != nil {
			log.Fatal("[2]", err)
		}

		finalOffset := decoder.InputOffset()
		obj := Obj{
			Data:  data,
			start: currentOffset,
			end:   finalOffset,
			bytes: finalOffset - currentOffset,
		}

		if index, exists := idMap[data.Id]; exists {
			log.Println("IDs should be unique. The id:", data.Id, " is already present and has value index of:", index)
			return &List{}, DuplicateIds
		}

		idMap[data.Id] = count
		*list = append(*list, obj)
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

	return list, None
}

func (s *Stream) Flush() {
	beforeBackupIDMap = map[string]uint64{}
	afterBackupIDMap = map[string]uint64{}
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
