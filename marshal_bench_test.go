package sb

import (
	"reflect"
	"testing"
)

func BenchmarkMarshalInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(42),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(float64(42)),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

type benchInt int

var _ SBMarshaler = benchInt(0)

func (b benchInt) MarshalSB(ctx Ctx, cont Proc) Proc {
	return ctx.Marshal(ctx, reflect.ValueOf(int(b)), cont)
}

func BenchmarkMarshalIntMarshaler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(benchInt(42)),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalArray(b *testing.B) {
	array := [8]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(array),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalStruct(b *testing.B) {
	s := struct {
		Foo int
		Bar int
	}{42, 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(s),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalMap(b *testing.B) {
	m := map[int]int{
		42: 42,
		1:  1,
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
	tuple := func() (int, int) {
		return 42, 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(tuple),
			Discard,
		); err != nil {
			b.Fatal(err)
		}
	}
}
