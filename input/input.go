package input

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Dalot/validate_backup/memstat"
)

type Input interface {
	InitReaders() (io.Reader, io.Reader)
}

type FilesInput struct {
	BeforeFileStr string
	AfterFileStr  string
}

func (i FilesInput) InitReaders() (io.Reader, io.Reader) {
	var beforeReader io.Reader
	var afterReader io.Reader
	beforeFile, err := os.Open(i.BeforeFileStr)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", i.BeforeFileStr, err.Error())
	}

	beforeFileInfo, err := beforeFile.Stat()
	if err != nil {
		log.Fatalf("Could not obtain stat, handle error: %v", err.Error())
	}

	afterFile, err := os.Open(i.AfterFileStr)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", i.AfterFileStr, err.Error())
	}

	afterFileInfo, err := afterFile.Stat()
	if err != nil {
		log.Fatalf("Could not obtain stat, handle error: %v", err.Error())
	}

	fmt.Printf("The [%s] is %s long\n", i.BeforeFileStr, memstat.FileSize(beforeFileInfo.Size()))
	fmt.Printf("The [%s] is %s long\n", i.AfterFileStr, memstat.FileSize(afterFileInfo.Size()))
	beforeReader = beforeFile
	afterReader = afterFile

	return beforeReader, afterReader
}

type BytesInput struct {
	BeforeFileStr string
	AfterFileStr  string
}

func (i BytesInput) InitReaders() (io.Reader, io.Reader) {
	var beforeReader io.Reader
	var afterReader io.Reader
	beforeBytesReader := bytes.NewReader([]byte(i.BeforeFileStr))
	afterBytesReader := bytes.NewReader([]byte(i.AfterFileStr))

	beforeReader = beforeBytesReader
	afterReader = afterBytesReader

	return beforeReader, afterReader
}
