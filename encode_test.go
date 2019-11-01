package sb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestEncodeBadStream(t *testing.T) {
	err := Encode(ioutil.Discard, NewDecoder(
		bytes.NewReader([]byte{
			KindString, // incomplete
		}),
	))
	if err == nil {
		t.Fatal()
	}
}

type badWriter struct{}

var _ io.Writer = badWriter{}

func (b badWriter) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

func TestEncodeToBadWriter(t *testing.T) {
	err := Encode(badWriter{}, NewMarshaler(42))
	if err == nil {
		t.Fatal()
	}
}
