package sb

import "testing"

func TestKindString(t *testing.T) {
	if KindInt.String() != "KindInt" {
		t.Fatal()
	}
}
