package sb

import "testing"

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

func BenchmarkMarshalStruct(b *testing.B) {
	type Foo struct {
		Foo int
		Bar string
		Baz []int
	}
	foo := Foo{
		Foo: 42,
		Bar: "42",
		Baz: []int{42, 42},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(foo)
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
