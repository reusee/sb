package sb

import (
	"reflect"
)

type Tuple []any

var _ SBMarshaler = Tuple{}

func (t Tuple) MarshalSB(
	vm ValueMarshalFunc,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindTuple,
		}, marshalTuple(vm, t, cont), nil
	}
}

func marshalTuple(vm ValueMarshalFunc, tuple Tuple, cont Proc) Proc {
	if len(tuple) == 0 {
		return func() (*Token, Proc, error) {
			return &Token{
				Kind: KindTupleEnd,
			}, cont, nil
		}
	}
	return vm(
		vm,
		reflect.ValueOf(tuple[0]),
		marshalTuple(
			vm,
			tuple[1:],
			cont,
		),
	)
}
