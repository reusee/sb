package sb

import (
	"bytes"
	"testing"
)

func TestUnmarshalArray(t *testing.T) {
	type foo [3]byte
	type S struct {
		Foos []foo
	}
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(S{
		Foos: []foo{
			foo{1},
			foo{2},
		},
	})); err != nil {
		t.Fatal(err)
	}
	var s S
	if err := Unmarshal(NewDecoder(buf), &s); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalNamedUint(t *testing.T) {
	type Foo uint32
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(Foo(42))); err != nil {
		t.Fatal(err)
	}
	var foo Foo
	if err := Unmarshal(NewDecoder(buf), &foo); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalStructWithPrivateField(t *testing.T) {
	type Foo struct {
		Bar int
		Foo int
	}
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(Foo{42, 42})); err != nil {
		t.Fatal(err)
	}
	type Bar struct {
		bar int
		Foo int
	}
	var bar Bar
	if err := Unmarshal(NewDecoder(buf), &bar); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalIncompleteStream(t *testing.T) {
	cases := [][]Token{
		{
			{Kind: KindObject},
		},
		{
			{Kind: KindObject},
			{KindString, "Foo"},
		},
		{
			{Kind: KindObject},
			{KindString, "Foo"},
			{KindString, "Bar"},
			{KindInt, 42},
		},
		{},
	}

	for _, c := range cases {
		var v any
		err := Unmarshal(List(c), &v)
		if err == nil {
			t.Fatal()
		}
	}

}
