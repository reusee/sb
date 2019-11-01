package sb

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestDecoderPeek(t *testing.T) {
	decoder := NewDecoder(bytes.NewReader([]byte{
		KindNaN,
	}))
	token, err := decoder.Peek()
	if err != nil {
		t.Fatal(err)
	}
	if token.Kind != KindNaN {
		t.Fatal()
	}
	token, err = decoder.Peek()
	if err != nil {
		t.Fatal(err)
	}
	if token.Kind != KindNaN {
		t.Fatal()
	}
}

type badReader struct{}

var _ io.Reader = badReader{}

func (r badReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

func TestDecodeBadReader(t *testing.T) {
	d := NewDecoder(badReader{})
	_, err := d.Next()
	if err == nil {
		t.Fatal()
	}
}
