package sb

import (
	"bytes"
	"testing"
)

func TestChunkedReader(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Copy(
		NewMarshaler(
			ChunkedReader{
				R: bytes.NewReader([]byte("foobarbaz")),
				N: 3,
			},
		),
		Encode(buf, nil),
	); err != nil {
		t.Fatal(err)
	}

	var data [][]byte
	if err := Unmarshal(
		NewDecoder(buf),
		&data,
	); err != nil {
		t.Fatal(err)
	}
	if len(data) != 3 {
		t.Fatal()
	}
	if !bytes.Equal(data[0], []byte("foo")) {
		t.Fatal()
	}
	if !bytes.Equal(data[1], []byte("bar")) {
		t.Fatal()
	}
	if !bytes.Equal(data[2], []byte("baz")) {
		t.Fatal()
	}
}
