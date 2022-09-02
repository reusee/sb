package sb

import (
	"fmt"
	"io"
	"reflect"

	"github.com/reusee/e5"
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
			return nil, we.With(WithPath(ctx), e5.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}
		if token.Kind != KindTuple {
			return nil, we.With(WithPath(ctx), e5.With(TypeMismatch(token.Kind, reflect.Func)))(UnmarshalError)
		}
		return t.unmarshal(0, ctx, cont), nil
	}
}

func (t *Tuple) unmarshal(i int, ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, we.With(WithPath(ctx), e5.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}
		if token.Kind == KindTupleEnd {
			return cont, nil
		}
		if i < len(*t) {
			var elem reflect.Value
			if (*t)[i] != nil {
				elem = reflect.New(reflect.TypeOf((*t)[i]))
				elem.Elem().Set(reflect.ValueOf((*t)[i]))
			} else {
				var value any
				elem = reflect.ValueOf(&value)
			}
			return ctx.Unmarshal(
				ctx.WithPath(i),
				elem,
				func(token *Token) (Sink, error) {
					(*t)[i] = elem.Elem().Interface()
					return t.unmarshal(i+1, ctx, cont)(token)
				},
			)(token)
		} else {
			var value any
			return ctx.Unmarshal(
				ctx.WithPath(i),
				reflect.ValueOf(&value),
				func(token *Token) (Sink, error) {
					*t = append(*t, value)
					return t.unmarshal(i+1, ctx, cont)(token)
				},
			)(token)
		}
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
			return nil, we.With(WithPath(ctx), e5.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}
		if token.Kind != KindTuple {
			return nil, we.With(WithPath(ctx), e5.With(TypeMismatch(token.Kind, reflect.Func)))(UnmarshalError)
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
			return nil, we.With(WithPath(ctx), e5.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}

		if token.Kind == KindTupleEnd {
			if i != len(types) {
				return nil, we.With(WithPath(ctx), e5.With(TooFewElement))(UnmarshalError)
			}
			return cont, nil
		}

		if i == len(types) {
			return nil, we.With(WithPath(ctx), e5.With(TooManyElement))(UnmarshalError)
		}

		if i < len(*target) {
			elem := reflect.New(types[i])
			if (*target)[i] != nil {
				elem.Elem().Set(reflect.ValueOf((*target)[i]))
			}
			return ctx.Unmarshal(
				ctx.WithPath(i),
				elem,
				func(token *Token) (Sink, error) {
					(*target)[i] = elem.Elem().Interface()
					i++
					return sink(token)
				},
			)(token)
		} else {
			elem := reflect.New(types[i])
			return ctx.Unmarshal(
				ctx.WithPath(i),
				elem,
				func(token *Token) (Sink, error) {
					*target = append(*target, elem.Elem().Interface())
					i++
					return sink(token)
				},
			)(token)
		}

	}

	return sink
}
