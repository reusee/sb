package sb

import (
	"bytes"
	"testing"
)

func TestDecodeJson(t *testing.T) {
	var n float64
	if err := Copy(
		DecodeJson(bytes.NewReader([]byte(`42`)), nil),
		Unmarshal(&n),
	); err != nil {
		t.Fatal(err)
	}
	if n != 42 {
		t.Fatal()
	}
}
