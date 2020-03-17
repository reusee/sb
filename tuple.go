package sb

import (
	"fmt"
	"reflect"
)

type Tuple []any

var _ SBMarshaler = Tuple{}

func (t Tuple) MarshalSB(ctx Ctx, cont Proc) Proc {
	return MarshalAsTuple(ctx, t, cont)
}

func MarshalAsTuple(ctx Ctx, tuple []any, cont Proc) Proc {
	var marshal Proc
	marshal = func() (*Token, Proc, error) {
		if len(tuple) == 0 {
			return &Token{
				Kind: KindTupleEnd,
			}, cont, nil
		}
		value := tuple[0]
		tuple = tuple[1:]
		return nil, ctx.Marshal(
			ctx,
			reflect.ValueOf(value),
			marshal,
		), nil
	}
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindTuple,
		}, marshal, nil
	}
}

var _ SBUnmarshaler = new(Tuple)

func (t *Tuple) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingTuple}
		}
		if token.Kind != KindTuple {
			return nil, UnmarshalError{ExpectingTuple}
		}
		return t.unmarshal(ctx, cont), nil
	}
}

func (t *Tuple) unmarshal(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingTuple}
		}
		if token.Kind == KindTupleEnd {
			return cont, nil
		}
		var value any
		return ctx.Unmarshal(
			ctx,
			reflect.ValueOf(&value),
			func(token *Token) (Sink, error) {
				*t = append(*t, value)
				return t.unmarshal(ctx, cont)(token)
			},
		)(token)
	}
}

func UnmarshalTupleTyped(ctx Ctx, typeSpec any, target *Tuple, cont Sink) Sink {
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

		return unmarshalTupleTyped(ctx, types, target, cont), nil
	}
}

func unmarshalTupleTyped(ctx Ctx, types []reflect.Type, target *Tuple, cont Sink) Sink {
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
		return ctx.Unmarshal(
			ctx,
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
