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
				byte(KindString), // incomplete
			}),
		),
		Encode(ioutil.Discard),
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
		Marshal(42),
		Encode(badWriter{}),
	); err == nil {
		t.Fatal()
	}
}

func BenchmarkEncodeInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(42),
			Encode(ioutil.Discard),
		); err != nil {
			b.Fatal(err)
		}
	}
}
