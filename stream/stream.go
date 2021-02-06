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
	Id   string
	Name string
}

type List map [string]Obj

type Stream struct{}

// Compare takes two readers, each to a json file, and compares them to find if they are equal
func Compare(beforeReader io.Reader, afterReader io.Reader) (*Stream, Comparison) {
	stream := &Stream{}

	listBeforeBackup, comp := stream.ParseBefore(beforeReader)
	if comp != None {
		return stream, comp
	}
	comp = stream.Compare(afterReader, listBeforeBackup)

	return stream, comp
}

// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) ParseBefore(input io.Reader) (map [string]Obj, Comparison)  {
	var count uint64 = 0
	var list = make(map [string]Obj)
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
			log.Fatal("[2]", err)
		}
		
		newObj := Obj{
			Name: obj.Name,
			Id: obj.Id,
		}
		list[obj.Id] = newObj

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


// Parse does stream processing of json file => less memory footprint, as we read record by record
func (s *Stream) Compare(input io.Reader, listBeforeBackup map[string]Obj) (Comparison)  {
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
			log.Fatal("[2]", err)
		}

		if existentObj, exists := listBeforeBackup[obj.Id]; !exists {
			log.Println("IDs should be unique. The id:", existentObj.Id, " is already present")
			return DuplicateIds
		} else {
			if existentObj.Name != obj.Name {
				log.Print("Files are different")
				log.Print("Before File: id(", existentObj.Id, "), value(", existentObj.Name,")")
				log.Print("After file: id(", obj.Id, "), value(", obj.Name,")")
				return Different

			} else {
				delete(listBeforeBackup, obj.Id)
			}
		}
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

	return Equal
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
