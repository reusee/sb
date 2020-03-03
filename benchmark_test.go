package sb

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(42),
			Discard,
		); err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkTreeFromMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			MustTreeFromStream(Marshal(42)).Iter(),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			NewPostHasher(Marshal(42), newMapHashState),
			Discard,
		); err != nil {
			b.Fatal(err)
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
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(benchFoo),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalStructAsTuple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(Tuple{
				benchFoo.Foo, benchFoo.Bar, benchFoo.Baz,
			}),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalStructAsTupleFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(func() (int, string, []int) {
				return benchFoo.Foo, benchFoo.Bar, benchFoo.Baz
			}),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTreeFromMarshalStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			MustTreeFromStream(Marshal(benchFoo)).Iter(),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			NewPostHasher(Marshal(benchFoo), newMapHashState),
			Discard,
		); err != nil {
			b.Fatal(err)
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
		if err := Copy(
			Marshal(m),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalTuple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(func() int {
				return 42
			}),
			Discard,
		); err != nil {
			b.Fatal(err)
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
