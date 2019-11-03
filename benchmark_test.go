package sb

import (
	"fmt"
	"testing"
)

func BenchmarkMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(42)
		for {
			token, err := m.Next()
			if err != nil {
				b.Fatal(err)
			}
			if token == nil {
				break
			}
		}
	}
}

func BenchmarkUnmarshalInt(b *testing.B) {
	tokens := MustTokensFromStream(NewMarshaler(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var n int
		if err := Unmarshal(tokens.Iter(), &n); err != nil {
			b.Fatal(err)
		}
		if n != 42 {
			b.Fatal()
		}
	}
}

type BenchFoo struct {
	Foo int
	Bar string
	Baz []int
}

var benchFoo = BenchFoo{
	Foo: 42,
	Bar: "42",
	Baz: []int{42, 42},
}

func BenchmarkMarshalStruct(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(benchFoo)
		for {
			token, err := m.Next()
			if err != nil {
				b.Fatal(err)
			}
			if token == nil {
				break
			}
		}
	}
}

func BenchmarkMarshalMap(b *testing.B) {
	m := make(map[int]string)
	for i := 0; i < 128; i++ {
		m[i] = fmt.Sprintf("%x", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(m)
		for {
			token, err := m.Next()
			if err != nil {
				b.Fatal(err)
			}
			if token == nil {
				break
			}
		}
	}
}
