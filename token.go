package sb

import (
	"fmt"
	"reflect"
)

type Kind uint8

const (
	KindInvalid Kind = iota

	KindArrayEnd
	KindObjectEnd

	KindBool
	KindInt
	KindUint
	KindFloat
	KindString
	KindIndirect
	KindNil

	KindArray
	KindObject
)

type Token struct {
	Kind  Kind
	Value any
}

func KindOf(kind reflect.Kind) Kind {
	switch kind {

	case reflect.Interface, reflect.Ptr:
		return KindIndirect

	case reflect.Bool:
		return KindBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return KindInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return KindUint

	case reflect.Float32, reflect.Float64:
		return KindFloat

	case reflect.Array, reflect.Slice:
		return KindArray

	case reflect.String:
		return KindString

	case reflect.Struct:
		return KindObject

	default:
		panic(fmt.Errorf("invalid kind: %s", kind.String()))

	}
}
