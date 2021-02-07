package stream

import (
	"bytes"
	"log"
	"testing"
)

var cases = []struct {
	a      string
	b      string
	result Comparison
	alternativeResult Comparison
}{
	/* 0 */{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Equal, None},
	{`[{"id": "7","name": "John doe"}]`, `[{"id": "7","name": "John doe"}]`, Equal, None},
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "42"}]`, Different, None},
	{`[{"id": "8","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "8","name": "John doe"}]`, Equal, None},
	{`[{"id": "4","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Different, None},
	/* 5 */{`[{"id": "5","name": "John doe"},{"id": "5","name": "John doe"}]`, `[{"id": "9","name": "John doe"},{"id": "5","name": "John doe"}]`, DuplicateIds, Different},
	{`[{"id": "5","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "6","name": "John doe"}]`, DuplicateIds, Different},
	{`[{"id": "6","name": "John doe"}]`, `[{}]`, Different, None},
	{`[{"id": "3","name": "John doe"}]`, `[{"id": "3"}]`, Different, None},
	{`[{"id": "3","name": "John doe"}]`, `[{"name": "John Doe}]`, Different, None},
	/* 10 */{`[{}]`, `[{}]`, Equal, None},
}

func TestCompare(t *testing.T) {
	for n := 1; n <= 10; n++ {
		for i, c := range cases {
			log.Println("CASE ", i)
			a := bytes.NewReader([]byte(c.a))
			b := bytes.NewReader([]byte(c.b))
	
			_, result := Compare(a, b)
	
			if result != c.result {
				if c.alternativeResult != None && result != c.alternativeResult {
					t.Errorf("case %d failed, got: %v, expected: %v", i, result, c.result)
				}
			}
			
		}
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

func getObjects(n int) map[string]Obj {
	list := make(map[string]Obj, n)
	for i := range list {
		list[i] = Obj{
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
