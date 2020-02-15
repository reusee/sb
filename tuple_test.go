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
