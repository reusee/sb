package sb

import (
	"bytes"
	"testing"
)

func TestUnmarshal(t *testing.T) {
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
