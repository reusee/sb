package sb

import (
	"bytes"
	"testing"
)

func TestDump(t *testing.T) {
	w := new(bytes.Buffer)
	dumpStream(NewMarshaler(42), "->", w)
	if w.String() != "->&{Kind:60 Value:42}\n" {
		t.Fatal()
	}
}
