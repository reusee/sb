package sb

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := Marshal(42)
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

func BenchmarkTreeFromMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := MustTreeFromStream(Marshal(42)).Iter()
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

func BenchmarkHashIntSha1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewPostHasher(Marshal(42), newMapHashState)
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
	tokens := MustTokensFromStream(Marshal(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var n int
		if err := Copy(tokens.Iter(), Unmarshal(&n)); err != nil {
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
		m := Marshal(benchFoo)
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

func BenchmarkTreeFromMarshalStruct(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := MustTreeFromStream(Marshal(benchFoo)).Iter()
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

func BenchmarkHashStructSha1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewPostHasher(Marshal(benchFoo), newMapHashState)
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
	tokens := MustTokensFromStream(Marshal(benchFoo))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var foo BenchFoo
		if err := Copy(tokens.Iter(), Unmarshal(&foo)); err != nil {
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
		m := Marshal(m)
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
		m := Marshal(func() int {
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
	tokens := MustTokensFromStream(Marshal(
		func() int {
			return 42
		},
	))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var fn func() int
		if err := Copy(tokens.Iter(), Unmarshal(&fn)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnarshalTupleCall(b *testing.B) {
	tokens := MustTokensFromStream(Marshal(
		func() int {
			return 42
		},
	))
	fn := func(int) {
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Copy(tokens.Iter(), Unmarshal(fn)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareTokens(b *testing.B) {
	tokens := MustTokensFromStream(
		Marshal(benchFoo),
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
			Marshal(benchFoo),
			Marshal(benchFoo),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareMarshalerAndTokens(b *testing.B) {
	tokens := MustTokensFromStream(
		Marshal(benchFoo),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			tokens.Iter(),
			Marshal(benchFoo),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecoders(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Copy(Marshal(benchFoo), Encode(buf)); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			Decode(bytes.NewReader(data)),
			Decode(bytes.NewReader(data)),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodersAndTokens(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Copy(Marshal(benchFoo), Encode(buf)); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	tokens := MustTokensFromStream(
		Marshal(benchFoo),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compare(
			tokens.Iter(),
			Decode(bytes.NewReader(data)),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
