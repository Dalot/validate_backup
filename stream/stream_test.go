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

		stream, result := Compare(a, b)

		if result != c.result {
			t.Errorf("case %d failed, got: %v, expected: %v", i, result, c.result)
		}

		stream.Flush()
	}
}

func BenchmarkSlice100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getSlice(100)
	}
}

func getSlice(n int) List {
	list := make([]*Obj, n)
	for i := range list {
		list[i] = &Obj{
			Name: "John",
			Id:   "Doe",
		}
	}
	return list
}
