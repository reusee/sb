package sb

import "testing"

func TestRef(t *testing.T) {
	value := []Ref{
		Ref("foo"),
		Ref("bar"),
		Ref("baz"),
	}
	var value2 []Ref
	if err := Copy(
		Marshal(value),
		Unmarshal(&value2),
	); err != nil {
		t.Fatal(err)
	}
	if res := MustCompare(
		Marshal(value),
		Marshal(value2),
	); res != 0 {
		t.Fatal()
	}
}
