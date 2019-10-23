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
	if err := Encode(buf, NewValue(S{
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
	if err := Encode(buf, NewValue(Foo(42))); err != nil {
		t.Fatal(err)
	}
	var foo Foo
	if err := Unmarshal(NewDecoder(buf), &foo); err != nil {
		t.Fatal(err)
	}
}
