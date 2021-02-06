package stream

import (
	"bytes"
	"testing"
)

var cases = []struct {
	a      string
	b      string
	result Comparison
}{
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Equal},
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Equal},
	{`[{"id": "5","name": "John doe"}]`, `[{"id": "5","name": "42"}]`, Different},
	{`[{"id": "5","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "5","name": "John doe"}]`, Equal},
	{`[{"id": "4","name": "John doe"}]`, `[{"id": "5","name": "John doe"}]`, Different},
	{`[{"id": "5","name": "John doe"},{"id": "5","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "5","name": "John doe"}]`, DuplicateIds},
	{`[{"id": "5","name": "John doe"},{"id": "6","name": "John doe"}]`, `[{"id": "6","name": "John doe"},{"id": "6","name": "John doe"}]`, DuplicateIds},
}

func TestCompare(t *testing.T) {
	for i, c := range cases {
		a := bytes.NewReader([]byte(c.a))
		b := bytes.NewReader([]byte(c.b))

		_, result := Compare(a, b)

		if result != c.result {
			t.Errorf("case %d failed, got: %v, expected: %v", i, result, c.result)
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

func getObjects(n int) List {
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
