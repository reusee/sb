package sb

import (
	"testing"
)

func TestConcat(t *testing.T) {
	var a any
	var b any
	if err := Copy(
		ConcatStreams(
			Marshal(42),
			Marshal(2),
		),
		ConcatSinks(
			Unmarshal(&a),
			Unmarshal(&b),
		),
	); err != nil {
		t.Fatal(err)
	}
	if a != 42 {
		t.Fatal()
	}
	if b != 2 {
		t.Fatal()
	}
}

func TestConcatNilSinks(t *testing.T) {
	if err := Copy(
		Marshal(42),
		ConcatSinks(),
	); err != nil {
		t.Fatal(err)
	}

	if err := Copy(
		Marshal(42),
		ConcatSinks(nil),
	); err != nil {
		t.Fatal(err)
	}

	var i int
	if err := Copy(
		Marshal(42),
		ConcatSinks(
			nil,
			nil,
			Unmarshal(&i),
			nil,
			nil,
		),
	); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal()
	}
}
