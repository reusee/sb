package sb

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	for _, v := range []any{
		true,
		false,
		int(42),
		uint(42),
		int32(42),
		uint64(42),
		"foo",
		strings.Repeat("foo", 1024),
		[]byte("foo"),
		[]byte(strings.Repeat("foo", 1024)),
	} {
		if err := Copy(
			Marshal(v),
			Encode(io.Discard),
		); err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		var l int
		if err := Copy(
			Marshal(v),
			Encode(buf),
			EncodedLen(&l, nil),
		); err != nil {
			t.Fatal(err)
		}
		if l != buf.Len() {
			t.Fatal()
		}
	}
}

func TestEncodeBadStream(t *testing.T) {
	if err := Copy(
		Decode(
			bytes.NewReader([]byte{
				byte(KindString), // incomplete
			}),
		),
		Encode(io.Discard),
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
			Encode(io.Discard),
		); err != nil {
			b.Fatal(err)
		}
	}
}
