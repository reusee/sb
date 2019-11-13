package sb

import (
	"bytes"
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

func BenchmarkUnmarshalStruct(b *testing.B) {
	tokens := MustTokensFromStream(NewMarshaler(benchFoo))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var foo BenchFoo
		if err := Unmarshal(tokens.Iter(), &foo); err != nil {
			b.Fatal(err)
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

func BenchmarkMarshalTuple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(func() int {
			return 42
		})
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

func BenchmarkUnarshalTuple(b *testing.B) {
	tokens := MustTokensFromStream(NewMarshaler(
		func() int {
			return 42
		},
	))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var fn func() int
		if err := Unmarshal(tokens.Iter(), &fn); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnarshalTupleCall(b *testing.B) {
	tokens := MustTokensFromStream(NewMarshaler(
		func() int {
			return 42
		},
	))
	fn := func(int) {
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(tokens.Iter(), fn); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareTokens(b *testing.B) {
	tokens := MustTokensFromStream(
		NewMarshaler(benchFoo),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			tokens.Iter(),
			tokens.Iter(),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareMarshaler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			NewMarshaler(benchFoo),
			NewMarshaler(benchFoo),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareMarshalerAndTokens(b *testing.B) {
	tokens := MustTokensFromStream(
		NewMarshaler(benchFoo),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			tokens.Iter(),
			NewMarshaler(benchFoo),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecoders(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(benchFoo)); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			NewDecoder(bytes.NewReader(data)),
			NewDecoder(bytes.NewReader(data)),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodersAndTokens(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(benchFoo)); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	tokens := MustTokensFromStream(
		NewMarshaler(benchFoo),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			tokens.Iter(),
			NewDecoder(bytes.NewReader(data)),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
