package stream

import (
	"log"
	"testing"

	"github.com/Dalot/validate_backup/input"
)

var cases = []struct {
	a                 string
	b                 string
	result            Comparison
	alternativeResult Comparison
}{
	/* 0 */
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Equal, None},
	{`[{"id": "7","name": "John doe"}]`, `[{"id": "7","name": "John doe"}]`, Equal, None},
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "42"}]`, Different, None},
	{`[{"id": "8","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "8","name": "John doe"}]`, Equal, None},
	{`[{"id": "4","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Different, None},
	/* 5 */
	{`[{"id": "5","name": "John doe"},{"id": "5","name": "John doe"}]`, `[{"id": "9","name": "John doe"},{"id": "5","name": "John doe"}]`, DuplicateIds, Different},
	{`[{"id": "5","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "6","name": "John doe"}]`, DuplicateIds, Different},
	{`[{"id": "6","name": "John doe"}]`, `[{}]`, InvalidJson, None},
	{`[{"id": "3","name": "John doe"}]`, `[{"id": "3"}]`, Different, None},
	{`[{"id": "3","name": "John doe"}]`, `[{"name": "John Doe}]`, InvalidJson, None},
	/* 10 */
	{`[{"id": "3","name": "John doe"}]`, `[{"name": "John Doe}]`, InvalidJson, None},
	{`[{"id": "3","name": "John doe", "job": "designer"}]`, `[{"id": "3","name": "John doe", "job": "designer"}]`, Equal, None},
	{`[{"id": 3,"name": "John doe", "job": "designer"}]`, `[{"id": 3,"name": "John doe", "job": "designer"}]`, Equal, None},
	{`[{"id": 3,"name": "John doe", "job": true }]`, `[{"id": 3,"name": "John doe", "job": true }]`, Equal, None},
	{`[{"id": 3,"name": "John doe", "job": true }]`, `[{"id": 3,"name": "John doe", "job": true }]`, Equal, None},
	// TODO: Support nested slices and objects of any type
	//{`[{"id": 3,"name": "John doe", "job": [false, true] }]`, `[{"id": 3,"name": "John doe", "job": [false, true] }]`, Equal, None},
	//{`[{"id": 3,"name": "John doe", "job": { "1": "great job", "2": "boring job"} }]`, `[{"id": 3,"name": "John doe", "job": [false, true] }]`, Equal, None},
	{`[{}]`, `[{}]`, InvalidJson, None},
}

func TestCompare(t *testing.T) {
	// run 10 times each test to reduce the probability of concurrency bugs
	for n := 1; n <= 10; n++ { 
		for i, c := range cases {
			log.Println("CASE ", i)
			bytesInput := input.BytesInput{
				BeforeFileStr: c.a,
				AfterFileStr:  c.b,
			}
			a, b := bytesInput.InitReaders()

			_, result := Compare(a, b)

			if result != c.result {
				hasAlternativeResult := c.alternativeResult != None
				if hasAlternativeResult {
					if result != c.alternativeResult {
						t.Errorf("case %d failed, got: %v, expected: %v , or %v", i, result, c.result, c.alternativeResult)
					}
				} else {
					t.Errorf("case %d failed, got: %v, expected: %v", i, result, c.result)
				}
			}

		}
	}
}

func TestCompareFiles(t *testing.T) {
	str1 := "test_files/before.json"
	str2 := "test_files/after.json"
	filesInput := input.FilesInput{
		BeforeFileStr: str1,
		AfterFileStr:  str2,
	}

	beforeReader, afterReader := filesInput.InitReaders()
	_, result := Compare(beforeReader, afterReader)
	if result != Equal {
		t.Errorf("case compare files failed , got: %v, expected: %v", result, Equal)
	}
}

func BenchmarkObjects(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getObjects(7000000)
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getSlice(7000000)
	}
}

func getObjects(n int) map[string]*Obj {
	list := make(map[string]*Obj, n)
	for i := range list {
		list[i] = &Obj{
			Name: "John",
			Id:   "Doe",
		}
	}
	return list
}

func getSlice(n int) []Obj {
	list := make([]Obj, n)
	for i := range list {
		list[i] = Obj{
			Name: "John",
			Id:   "Doe",
		}
	}
	return list
}
