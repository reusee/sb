package sb

import (
	"encoding"
	"reflect"
	"sort"

	"github.com/reusee/e4"
)

type SBMarshaler interface {
	MarshalSB(ctx Ctx, cont Proc) Proc
}

func Marshal(value any) *Proc {
	marshaler := MarshalValue(Ctx{
		Marshal: MarshalValue,
	}, reflect.ValueOf(value), nil)
	return &marshaler
}

func MarshalCtx(ctx Ctx, value any) *Proc {
	marshaler := MarshalValue(ctx, reflect.ValueOf(value), nil)
	return &marshaler
}

func TapMarshal(ctx Ctx, value any, fn func(Ctx, reflect.Value)) *Proc {
	marshal := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		fn(ctx, value)
		return MarshalValue(ctx, value, cont)
	}
	ctx.Marshal = marshal
	proc := marshal(ctx, reflect.ValueOf(value), nil)
	return &proc
}

type UnmarshalFunc func(Ctx, Sink) Sink

var _ SBUnmarshaler = UnmarshalFunc(nil)

func (f UnmarshalFunc) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return f(ctx, cont)
}

func MarshalValue(ctx Ctx, value reflect.Value, cont Proc) Proc {
	if ctx.Marshal == nil {
		ctx.Marshal = MarshalValue
	}

	marshal := func() (*Token, Proc, error) {

		if value.IsValid() {

			switch v := value.Interface().(type) {

			case int:
				return &Token{Kind: KindInt, Value: v}, cont, nil

			case string:
				return &Token{Kind: KindString, Value: v}, cont, nil

			case bool:
				return &Token{Kind: KindBool, Value: v}, cont, nil

			case uint32:
				return &Token{Kind: KindUint32, Value: v}, cont, nil

			case int64:
				return &Token{Kind: KindInt64, Value: v}, cont, nil

			case uint64:
				return &Token{Kind: KindUint64, Value: v}, cont, nil

			case uint16:
				return &Token{Kind: KindUint16, Value: v}, cont, nil

			case uint8:
				return &Token{Kind: KindUint8, Value: v}, cont, nil

			case int32:
				return &Token{Kind: KindInt32, Value: v}, cont, nil

			case uint:
				return &Token{Kind: KindUint, Value: v}, cont, nil

			case float64:
				if v != v {
					return NaN, cont, nil
				} else {
					return &Token{Kind: KindFloat64, Value: v}, cont, nil
				}

			case int8:
				return &Token{Kind: KindInt8, Value: v}, cont, nil

			case float32:
				if v != v {
					return NaN, cont, nil
				} else {
					return &Token{Kind: KindFloat32, Value: v}, cont, nil
				}

			case int16:
				return &Token{Kind: KindInt16, Value: v}, cont, nil

			case SBMarshaler:
				if value.Kind() == reflect.Ptr && value.IsNil() {
					_, found := value.Type().Elem().MethodByName("SBMarshaler")
					if !found {
						// SBMarshaler is method defined on non-pointer type
						// calling SBMarshaler will deref the nil pointer
						// do not try to
						return Nil, cont, nil
					}
				}
				return nil, v.MarshalSB(ctx, cont), nil

			case encoding.BinaryMarshaler:
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, nil, we(MarshalError, WithPath(ctx), e4.With(err))
				}
				return nil, ctx.Marshal(ctx, reflect.ValueOf(string(bs)), cont), nil

			case encoding.TextMarshaler:
				bs, err := v.MarshalText()
				if err != nil {
					return nil, nil, we(MarshalError, WithPath(ctx), e4.With(err))
				}
				return nil, ctx.Marshal(ctx, reflect.ValueOf(string(bs)), cont), nil

			case *Token:
				return v, cont, nil

			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			return Nil, cont, nil

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				return Nil, cont, nil
			} else {
				ctx.pointerDepth++
				if ctx.pointerDepth == 1000 {
					ctx.detectCycleEnabled = true
				}
				if ctx.detectCycleEnabled {
					ptr := value.Pointer()
					for _, p := range ctx.visitedPointers {
						if p == ptr {
							return nil, nil, we(MarshalError, WithPath(ctx), e4.With(CyclicPointer))
						}
					}
					ctx.visitedPointers = append(ctx.visitedPointers, ptr)
				}
				return nil, ctx.Marshal(ctx, value.Elem(), cont), nil
			}

		case reflect.Bool:
			return &Token{
				Kind:  KindBool,
				Value: bool(value.Bool()),
			}, cont, nil

		case reflect.Int:
			return &Token{
				Kind:  KindInt,
				Value: int(value.Int()),
			}, cont, nil

		case reflect.Int8:
			return &Token{
				Kind:  KindInt8,
				Value: int8(value.Int()),
			}, cont, nil

		case reflect.Int16:
			return &Token{
				Kind:  KindInt16,
				Value: int16(value.Int()),
			}, cont, nil

		case reflect.Int32:
			return &Token{
				Kind:  KindInt32,
				Value: int32(value.Int()),
			}, cont, nil

		case reflect.Int64:
			return &Token{
				Kind:  KindInt64,
				Value: int64(value.Int()),
			}, cont, nil

		case reflect.Uint:
			return &Token{
				Kind:  KindUint,
				Value: uint(value.Uint()),
			}, cont, nil

		case reflect.Uint8:
			return &Token{
				Kind:  KindUint8,
				Value: uint8(value.Uint()),
			}, cont, nil

		case reflect.Uint16:
			return &Token{
				Kind:  KindUint16,
				Value: uint16(value.Uint()),
			}, cont, nil

		case reflect.Uint32:
			return &Token{
				Kind:  KindUint32,
				Value: uint32(value.Uint()),
			}, cont, nil

		case reflect.Uint64:
			return &Token{
				Kind:  KindUint64,
				Value: uint64(value.Uint()),
			}, cont, nil

		case reflect.Float32:
			f := value.Float()
			if f != f {
				return NaN, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat32,
					Value: float32(f),
				}, cont, nil
			}

		case reflect.Float64:
			f := value.Float()
			if f != f {
				return NaN, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat64,
					Value: f,
				}, cont, nil
			}

		case reflect.Array, reflect.Slice:
			if isBytes(value.Type()) {
				return &Token{
					Kind:  KindBytes,
					Value: toBytes(value),
				}, cont, nil
			} else {
				return nil, MarshalArray(ctx, value, 0, cont), nil
			}

		case reflect.String:
			return &Token{
				Kind:  KindString,
				Value: value.String(),
			}, cont, nil

		case reflect.Struct:
			return nil, MarshalStruct(ctx, value, cont), nil

		case reflect.Map:
			return nil, MarshalMap(ctx, value, cont), nil

		case reflect.Func:
			if value.Type().NumIn() != 0 {
				return nil, nil, we(
					MarshalError,
					WithPath(ctx),
					e4.With(BadTupleType),
					e4.NewInfo("bad tuple type: %v", value.Type()),
				)
			}
			items := value.Call([]reflect.Value{})
			return nil, MarshalTuple(
				ctx,
				items,
				cont,
			), nil

		default:
			return nil, cont, nil

		}
	}

	if value.IsValid() {
		if name, ok := registeredTypeToName.Load(value.Type()); ok {
			return func() (*Token, Proc, error) {
				return &Token{
					Kind:  KindTypeName,
					Value: name.(string),
				}, marshal, nil
			}
		}
	}

	return marshal
}

var arrayEndToken = reflect.ValueOf(&Token{
	Kind: KindArrayEnd,
})

var arrayToken = &Token{
	Kind: KindArray,
}

func MarshalArray(ctx Ctx, value reflect.Value, index int, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if index >= value.Len() {
			return nil, ctx.Marshal(
				ctx,
				arrayEndToken,
				cont,
			), nil
		}
		v := value.Index(index)
		index++
		return nil, ctx.Marshal(
			ctx.WithPath(index-1),
			v,
			proc,
		), nil
	}
	return func() (*Token, Proc, error) {
		return arrayToken, proc, nil
	}
}

var objectEndToken = reflect.ValueOf(&Token{
	Kind: KindObjectEnd,
})

var objectToken = &Token{
	Kind: KindObject,
}

func MarshalStruct(ctx Ctx, value reflect.Value, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return objectToken, MarshalStructFields(ctx, value, cont), nil
	}
}

func MarshalStructFields(ctx Ctx, value reflect.Value, cont Proc) Proc {
	valueType := value.Type()
	numField := valueType.NumField()
	fieldIdx := 0
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if fieldIdx == numField {
			return nil, ctx.Marshal(
				ctx,
				objectEndToken,
				cont,
			), nil
		}

		field := valueType.Field(fieldIdx)
		if ctx.SkipEmptyStructFields {
			fieldValue := value.Field(fieldIdx)
			if fieldValue.IsZero() {
				fieldIdx++
				return nil, proc, nil
			}
			if field.Type.Kind() == reflect.Slice && fieldValue.Len() == 0 {
				fieldIdx++
				return nil, proc, nil
			}
		}
		if field.PkgPath != "" {
			// unexported field
			fieldIdx++
			return nil, proc, nil
		}

		fieldIdx++
		return nil, ctx.Marshal(
			ctx.WithPath(field.Name),
			reflect.ValueOf(field.Name),
			func() (*Token, Proc, error) {
				return nil, ctx.Marshal(
					ctx.WithPath(field.Name),
					value.FieldByIndex(field.Index),
					proc,
				), nil
			},
		), nil
	}
	return proc
}

type MapTuple struct {
	Key       reflect.Value
	Value     reflect.Value
	KeyTokens Tokens
}

func MarshalMap(ctx Ctx, value reflect.Value, cont Proc) Proc {
	return MarshalMapIter(
		ctx,
		value,
		value.MapRange(),
		make([]*MapTuple, 0, value.Len()),
		cont,
	)
}

var mapEndToken = reflect.ValueOf(&Token{
	Kind: KindMapEnd,
})

var mapToken = &Token{
	Kind: KindMap,
}

func MarshalMapIter(ctx Ctx, value reflect.Value, iter *reflect.MapIter, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if !iter.Next() {
			// done
			sort.Slice(tuples, func(i, j int) bool {
				return MustCompare(
					tuples[i].KeyTokens.Iter(),
					tuples[j].KeyTokens.Iter(),
				) < 0
			})
			return nil, MarshalMapTuples(ctx, tuples, cont), nil
		}
		var tokens Tokens
		keyMarshalProc := MarshalValue(Ctx{}, iter.Key(), nil)
		if err := Copy(
			// tokens are for sorting only, so do not call ctx.Marshal
			&keyMarshalProc,
			CollectTokens(&tokens),
		); err != nil {
			return nil, nil, we(MarshalError, WithPath(ctx), e4.With(err))
		}
		if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, nil, we(MarshalError, WithPath(ctx), e4.With(BadMapKey))
		}
		tuples = append(tuples, &MapTuple{
			KeyTokens: tokens,
			Key:       iter.Key(),
			Value:     iter.Value(),
		})
		return nil, proc, nil
	}
	return func() (*Token, Proc, error) {
		return mapToken, proc, nil
	}
}

func MarshalMapTuples(ctx Ctx, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(tuples) == 0 {
			return nil, ctx.Marshal(
				ctx,
				mapEndToken,
				cont,
			), nil
		}
		tuple := tuples[0]
		tuples = tuples[1:]
		path := tuple.Key.Interface()
		return nil, ctx.Marshal(
			ctx.WithPath(path),
			tuple.Key,
			// must wrap in closure to delay value marshaling
			func() (*Token, Proc, error) {
				return nil, ctx.Marshal(
					ctx.WithPath(path),
					tuple.Value,
					proc,
				), nil
			},
		), nil
	}
	return proc
}

var tupleEndToken = reflect.ValueOf(&Token{
	Kind: KindTupleEnd,
})

var tupleToken = &Token{
	Kind: KindTuple,
}

func MarshalTuple(ctx Ctx, items []reflect.Value, cont Proc) Proc {
	var proc Proc
	var n int
	proc = func() (*Token, Proc, error) {
		if len(items) == 0 {
			return nil, ctx.Marshal(
				ctx,
				tupleEndToken,
				cont,
			), nil
		} else {
			v := items[0]
			items = items[1:]
			n++
			return nil, ctx.Marshal(
				ctx.WithPath(n-1),
				v,
				proc,
			), nil
		}
	}
	return func() (*Token, Proc, error) {
		return tupleToken, proc, nil
	}
}
