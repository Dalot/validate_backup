package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Dalot/validate_backup/memstat"
	"github.com/Dalot/validate_backup/stream"
)

func main() {
	beforeBackupPtr := flag.String("f1", "", "File path before backup.")
	afterBackupPtr := flag.String("f2", "", "File path after backup.")
	flag.Parse()
	start := time.Now()

	//str1 := "xl_before.json"
	//str1 := "before.json"
	//str2 := "xl_after.json"
	//str2 := "after.json"
	//beforeBackupPtr = &str1
	//afterBackupPtr = &str2

	log.Println("before: ", *beforeBackupPtr, "after: ", *afterBackupPtr)
	beforeReader, afterReader := reader.InitReaders(*beforeBackupPtr, *afterBackupPtr)
	// this program is about parsing a large json file using a small memory footprint.
	// you may generate data using generate_data.sh it generate ~900mB in ~600 seconds.
	go memstat.PrintUsage()
	_, result := stream.Compare(beforeReader, afterReader)

	fmt.Println("RESULT: ", result.String())
	elapsed := time.Since(start)

	fmt.Printf("To parse the file took [%v]\n", elapsed)

}
