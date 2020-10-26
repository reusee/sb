package sb

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type badReader struct{}

var _ io.Reader = badReader{}

func (r badReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

type badByteReader struct{}

var _ io.Reader = badByteReader{}

func (r badByteReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("bad")
}

func (r badByteReader) ReadByte() (byte, error) {
	return 0, fmt.Errorf("bad")
}

func TestDecodeBadReader(t *testing.T) {
	d := Decode(badReader{})
	_, err := d.Next()
	if err == nil {
		t.Fatal()
	}
}

func TestDecodeBadByteReader(t *testing.T) {
	d := Decode(badByteReader{})
	_, err := d.Next()
	if err == nil {
		t.Fatal()
	}
}

func TestDecodeError(t *testing.T) {
	for _, kind := range []Kind{
		KindBool,
		KindInt, KindInt8, KindInt16, KindInt32, KindInt64,
		KindUint, KindUint8, KindUint16, KindUint32, KindUint64,
		KindFloat32, KindFloat64,
		KindString, KindBytes,
	} {
		err := Copy(
			Decode(bytes.NewReader([]byte{byte(kind)})),
			Discard,
		)
		if err == nil {
			t.Fatal()
		}
		if err.Error() != "DecodeError: EOF at 1" {
			t.Fatal()
		}
	}

}
