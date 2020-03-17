package sb

import (
	"bytes"
	"reflect"
	"testing"
)

func TestToBytes(t *testing.T) {
	if !bytes.Equal(
		toBytes(reflect.ValueOf([3]byte{1, 1, 1})),
		[]byte{1, 1, 1},
	) {
		t.Fatal()
	}
}
