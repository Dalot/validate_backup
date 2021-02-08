package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Dalot/validate_backup/input"
	"github.com/Dalot/validate_backup/memstat"
	"github.com/Dalot/validate_backup/stream"
)

func main() {
	beforeBackupPtr := flag.String("f1", "", "File path before backup.")
	afterBackupPtr := flag.String("f2", "", "File path after backup.")
	flag.Parse()
	start := time.Now()

	log.Println("before: ", *beforeBackupPtr, "after: ", *afterBackupPtr)
	filesInput := input.FilesInput{
		BeforeFileStr: *beforeBackupPtr,
		AfterFileStr:  *afterBackupPtr,
	}
	beforeReader, afterReader := filesInput.InitReaders()

	go memstat.PrintUsage()
	_, result := stream.Compare(beforeReader, afterReader)

	fmt.Println("RESULT: ", result.String())
	elapsed := time.Since(start)

	fmt.Printf("To parse the file took [%v]\n", elapsed)

}
