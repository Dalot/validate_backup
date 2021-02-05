package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Dalot/validate_backup/memstat"
	"github.com/Dalot/validate_backup/stream"
)

func main() {
	beforeBackupPtr := flag.String("f1", "", "File path before backup.")
	afterBackupPtr := flag.String("f2", "", "File path after backup.")
	flag.Parse()
	start := time.Now()

	log.Println("before: ", *beforeBackupPtr, "after: ", *afterBackupPtr)
	beforeReader, afterReader := initReaders(*beforeBackupPtr, *afterBackupPtr)
	// this program is about parsing a large json file using a small memory footprint.
	// you may generate data using generate_data.sh it generate ~900mB in ~600 seconds.
	go memstat.PrintUsage()
	_, result := stream.Compare(beforeReader, afterReader)

	fmt.Println("RESULT: ", result.String())
	elapsed := time.Since(start)

	fmt.Printf("To parse the file took [%v]\n", elapsed)

}

func initReaders(beforeFileStr string, afterFileStr string) (*bufio.Reader, *bufio.Reader) {
	beforeFile, err := os.Open(beforeFileStr)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", beforeFileStr, err.Error())
	}

	beforeFileInfo, err := beforeFile.Stat()
	if err != nil {
		log.Fatalf("Could not obtain stat, handle error: %v", err.Error())
	}

	beforeReader := bufio.NewReader(beforeFile)

	afterFile, err := os.Open(afterFileStr)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", afterFileStr, err.Error())
	}

	afterFileInfo, err := afterFile.Stat()
	if err != nil {
		log.Fatalf("Could not obtain stat, handle error: %v", err.Error())
	}

	afterReader := bufio.NewReader(afterFile)

	fmt.Printf("The [%s] is %s long\n", beforeFileStr, memstat.FileSize(beforeFileInfo.Size()))
	fmt.Printf("The [%s] is %s long\n", afterFileStr, memstat.FileSize(afterFileInfo.Size()))

	return beforeReader, afterReader
}
