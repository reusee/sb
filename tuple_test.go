package sb

import "testing"

func TestTuple(t *testing.T) {
	s := Marshal(Tuple{
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
	if err := Copy(
		Marshal(Tuple{
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

	if err := Copy(
		Marshal(Tuple{
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

	if err := Copy(
		Tokens{}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingTuple) {
		t.Fatal(err)
	}

	if err := Copy(
		Tokens{
			Token{
				Kind: KindString,
			},
		}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingTuple) {
		t.Fatal(err)
	}

	if err := Copy(
		Tokens{
			Token{
				Kind: KindTuple,
			},
		}.Iter(),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingValue) {
		t.Fatal(err)
	}

	if err := Copy(
		Marshal(Tuple{
			42, true, "foo",
		}),
		UnmarshalTupleTyped(func(int, bool) {}, &tuple, nil),
	); !is(err, TooManyElement) {
		t.Fatal(err)
	}

	if err := Copy(
		Marshal(Tuple{
			42, true,
		}),
		UnmarshalTupleTyped(func(int, bool, string) {}, &tuple, nil),
	); !is(err, ExpectingValue) {
		t.Fatal(err)
	}

}
