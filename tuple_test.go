package sb

import "testing"

func TestTuple(t *testing.T) {
	s := NewMarshaler(Tuple{
		42, true, "foo",
	})
	if err := Unmarshal(s, func(i int, b bool, s string) {
		if i != 42 {
			t.Fatal()
		}
		if !b {
			t.Fatal()
		}
		if s != "foo" {
			t.Fatal()
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestTupleUnmarshalTyped(t *testing.T) {

	var tuple Tuple
	if err := Pipe(
		NewMarshaler(Tuple{
			42, true, "foo",
		}),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); err != nil {
		t.Fatal(err)
	}
	if tuple[0] != 42 {
		t.Fatal()
	}
	if tuple[1] != true {
		t.Fatal()
	}
	if tuple[2] != "foo" {
		t.Fatal()
	}

	if err := Pipe(
		NewMarshaler(Tuple{
			1, false, "bar",
		}),
		UnmarshalTupleTyped(struct {
			int
			bool
			string
		}{}, &tuple, nil),
	); err != nil {
		t.Fatal(err)
	}
	if tuple[0] != 42 {
		t.Fatal()
	}
	if tuple[1] != true {
		t.Fatal()
	}
	if tuple[2] != "foo" {
		t.Fatal()
	}

	if err := Pipe(
		Tokens{}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingTuple) {
		t.Fatal(err)
	}

	if err := Pipe(
		Tokens{
			Token{
				Kind: KindString,
			},
		}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingTuple) {
		t.Fatal(err)
	}

	if err := Pipe(
		Tokens{
			Token{
				Kind: KindTuple,
			},
		}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingValue) {
		t.Fatal(err)
	}

	if err := Pipe(
		NewMarshaler(Tuple{
			42, true, "foo",
		}),
		UnmarshalTupleTyped(func(int, bool) {}, &tuple, nil),
	); !is(err, TooManyElement) {
		t.Fatal(err)
	}

	if err := Pipe(
		NewMarshaler(Tuple{
			42, true,
		}),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingValue) {
		t.Fatal(err)
	}

}
