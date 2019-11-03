package sb

import (
	"reflect"
	"testing"
)

func TestToBytes(t *testing.T) {
	toBytes(reflect.ValueOf([3]byte{1}))
}
