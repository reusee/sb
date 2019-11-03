package sb

import (
	"fmt"
	"reflect"
)

var (
	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
	byteType  = reflect.TypeOf((*byte)(nil)).Elem()
)

func isBytes(t reflect.Type) bool {
	if t.AssignableTo(bytesType) {
		return true
	}
	if t.Kind() == reflect.Array && t.Elem() == byteType {
		return true
	}
	return false
}

func toBytes(v reflect.Value) []byte {
	t := v.Type()
	if t == bytesType {
		return v.Interface().([]byte)
	} else if t.AssignableTo(bytesType) {
		var bs []byte
		reflect.ValueOf(&bs).Elem().Set(v)
		return bs
	} else if t.Kind() == reflect.Array {
		return toBytes(v.Slice(0, v.Len()))
	}
	panic(fmt.Errorf("unknown type %v", t)) // NOCOVER
}
