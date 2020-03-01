package sb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestEncodeBadStream(t *testing.T) {
	if err := Copy(
		Decode(
			bytes.NewReader([]byte{
				KindString, // incomplete
			}),
		),
		Encode(ioutil.Discard, nil),
	); err == nil {
		t.Fatal()
	}
}

type badWriter struct{}

var _ io.Writer = badWriter{}

func (b badWriter) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

func TestEncodeToBadWriter(t *testing.T) {
	if err := Copy(
		NewMarshaler(42),
		Encode(badWriter{}, nil),
	); err == nil {
		t.Fatal()
	}
}
