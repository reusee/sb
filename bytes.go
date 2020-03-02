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
	if t.Kind() == reflect.Array && t.Elem().AssignableTo(byteType) {
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
		bs := make([]byte, v.Len())
		reflect.Copy(reflect.ValueOf(bs), v)
		return bs
	}
	panic(fmt.Errorf("non-bytes type %v", t)) // NOCOVER
}
