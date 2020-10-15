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
			return nil, NewUnmarshalError(ctx, ExpectingTuple)
		}
		if token.Kind != KindTuple {
			return nil, NewUnmarshalError(ctx, ExpectingTuple)
		}
		return t.unmarshal(ctx, cont), nil
	}
}

func (t *Tuple) unmarshal(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, NewUnmarshalError(ctx, ExpectingValue)
		}
		if token.Kind == KindTupleEnd {
			return cont, nil
		}
		var value any
		return ctx.Unmarshal(
			ctx.WithPath(len(*t)),
			reflect.ValueOf(&value),
			func(token *Token) (Sink, error) {
				*t = append(*t, value)
				return t.unmarshal(ctx, cont)(token)
			},
		)(token)
	}
}

func TupleTypes(typeSpec any) []reflect.Type {
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
	return types
}

func UnmarshalTupleTyped(ctx Ctx, types []reflect.Type, target *Tuple, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, NewUnmarshalError(ctx, ExpectingTuple)
		}
		if token.Kind != KindTuple {
			return nil, NewUnmarshalError(ctx, ExpectingTuple)
		}
		return unmarshalTupleTyped(ctx, types, target, cont), nil
	}
}

type TypedTuple struct {
	Types  []reflect.Type
	Values Tuple
}

var _ SBUnmarshaler = new(TypedTuple)

func (t *TypedTuple) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return UnmarshalTupleTyped(ctx, t.Types, &t.Values, cont)
}

func unmarshalTupleTyped(ctx Ctx, types []reflect.Type, target *Tuple, cont Sink) Sink {
	var sink Sink
	i := 0
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return nil, NewUnmarshalError(ctx, ExpectingValue)
		}

		if token.Kind == KindTupleEnd {
			if i != len(types) {
				return nil, NewUnmarshalError(ctx, ExpectingValue)
			}
			return cont, nil
		}

		if i == len(types) {
			return nil, NewUnmarshalError(ctx, TooManyElement)
		}

		var elem reflect.Value
		if i < len(*target) {
			elem = reflect.ValueOf(*target).Index(i).Addr()
		} else {
			elem = reflect.New(types[i])
		}
		return ctx.Unmarshal(
			ctx.WithPath(len(*target)),
			elem,
			func(token *Token) (Sink, error) {
				if i < len(*target) {
					(*target)[i] = elem.Elem().Interface()
				} else {
					*target = append(*target, elem.Elem().Interface())
				}
				i++
				return sink(token)
			},
		)(token)

	}

	return sink
}
