package sb

import (
	"fmt"
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

func UnmarshalTupleTyped(typeSpec any, target *Tuple, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingTuple}
		}
		if token.Kind != KindTuple {
			return nil, UnmarshalError{ExpectingTuple}
		}

		var types []reflect.Type
		spec := reflect.ValueOf(typeSpec)
		specType := spec.Type()
		switch specType.Kind() {
		case reflect.Func:
			for i := 0; i < specType.NumIn(); i++ {
				types = append(types, specType.In(i))
			}
		case reflect.Struct:
			for i := 0; i < specType.NumField(); i++ {
				types = append(types, specType.Field(i).Type)
			}
		default: // NOCOVER
			panic(fmt.Errorf("bad type: %T", typeSpec))
		}

		return unmarshalTupleTyped(types, target, cont), nil
	}
}

func unmarshalTupleTyped(types []reflect.Type, target *Tuple, cont Sink) Sink {
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingValue}
		}

		if token.Kind == KindTupleEnd {
			if len(types) > 0 {
				return nil, UnmarshalError{ExpectingValue}
			}
			return cont, nil
		}

		if len(types) == 0 {
			return nil, UnmarshalError{TooManyElement}
		}

		elem := reflect.New(types[0])
		return UnmarshalValue(
			elem,
			func(token *Token) (Sink, error) {
				*target = append(*target, elem.Elem().Interface())
				types = types[1:]
				return sink(token)
			},
		)(token)

	}

	return sink
}
