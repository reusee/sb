package sb

import (
	"fmt"
	"io"
	"testing"
)

type badReader struct{}

var _ io.Reader = badReader{}

func (r badReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

func TestDecodeBadReader(t *testing.T) {
	d := Decode(badReader{})
	_, err := d.Next()
	if err == nil {
		t.Fatal()
	}
}
