package sb

import (
	"bytes"
	"testing"
)

func TestDump(t *testing.T) {
	w := new(bytes.Buffer)
	dumpStream(Marshal(42), "->", w)
	if w.String() != "->&{Kind:KindInt Value:42}\n" {
		t.Fatalf("got %q", w.String())
	}
}
