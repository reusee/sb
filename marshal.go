package sb

import (
	"encoding"
	"reflect"

	"github.com/reusee/e5"
	"slices"
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

func MarshalValue(ctx Ctx, value reflect.Value, cont Proc) Proc {
	if ctx.Marshal == nil {
		ctx.Marshal = MarshalValue
	}

	marshal := func(token *Token) (Proc, error) {

		if value.IsValid() {

			switch v := value.Interface().(type) {

			case int:
				token.Kind = KindInt
				token.Value = v
				return cont, nil

			case string:
				token.Kind = KindString
				token.Value = v
				return cont, nil

			case bool:
				token.Kind = KindBool
				token.Value = v
				return cont, nil

			case uint32:
				token.Kind = KindUint32
				token.Value = v
				return cont, nil

			case int64:
				token.Kind = KindInt64
				token.Value = v
				return cont, nil

			case uint64:
				token.Kind = KindUint64
				token.Value = v
				return cont, nil

			case uintptr:
				token.Kind = KindPointer
				token.Value = v
				return cont, nil

			case uint16:
				token.Kind = KindUint16
				token.Value = v
				return cont, nil

			case uint8:
				token.Kind = KindUint8
				token.Value = v
				return cont, nil

			case int32:
				token.Kind = KindInt32
				token.Value = v
				return cont, nil

			case uint:
				token.Kind = KindUint
				token.Value = v
				return cont, nil

			case float64:
				if v != v {
					*token = NaN
					return cont, nil
				} else {
					token.Kind = KindFloat64
					token.Value = v
					return cont, nil
				}

			case int8:
				token.Kind = KindInt8
				token.Value = v
				return cont, nil

			case float32:
				if v != v {
					*token = NaN
					return cont, nil
				} else {
					token.Kind = KindFloat32
					token.Value = v
					return cont, nil
				}

			case int16:
				token.Kind = KindInt16
				token.Value = v
				return cont, nil

			case SBMarshaler:
				if value.Kind() == reflect.Ptr && value.IsNil() {
					_, found := value.Type().Elem().MethodByName("SBMarshaler")
					if !found {
						// SBMarshaler is method defined on non-pointer type
						// calling SBMarshaler will deref the nil pointer
						// do not try to
						*token = Nil
						return cont, nil
					}
				}
				return v.MarshalSB(ctx, cont), nil

			case encoding.BinaryMarshaler:
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, we.With(e5.With(MarshalError), WithPath(ctx))(err)
				}
				return ctx.Marshal(ctx, reflect.ValueOf(string(bs)), cont), nil

			case encoding.TextMarshaler:
				bs, err := v.MarshalText()
				if err != nil {
					return nil, we.With(e5.With(MarshalError), WithPath(ctx))(err)
				}
				return ctx.Marshal(ctx, reflect.ValueOf(string(bs)), cont), nil

			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			*token = Nil
			return cont, nil

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				*token = Nil
				return cont, nil
			} else {
				ctx.pointerDepth++
				if ctx.pointerDepth == 1000 {
					ctx.detectCycleEnabled = true
				}
				if ctx.detectCycleEnabled {
					ptr := value.Pointer()
					for _, p := range ctx.visitedPointers {
						if p == ptr {
							return nil, we.With(WithPath(ctx), e5.With(CyclicPointer))(MarshalError)
						}
					}
					ctx.visitedPointers = append(ctx.visitedPointers, ptr)
				}
				return ctx.Marshal(ctx, value.Elem(), cont), nil
			}

		case reflect.Bool:
			token.Kind = KindBool
			token.Value = bool(value.Bool())
			return cont, nil

		case reflect.Int:
			token.Kind = KindInt
			token.Value = int(value.Int())
			return cont, nil

		case reflect.Int8:
			token.Kind = KindInt8
			token.Value = int8(value.Int())
			return cont, nil

		case reflect.Int16:
			token.Kind = KindInt16
			token.Value = int16(value.Int())
			return cont, nil

		case reflect.Int32:
			token.Kind = KindInt32
			token.Value = int32(value.Int())
			return cont, nil

		case reflect.Int64:
			token.Kind = KindInt64
			token.Value = int64(value.Int())
			return cont, nil

		case reflect.Uint:
			token.Kind = KindUint
			token.Value = uint(value.Uint())
			return cont, nil

		case reflect.Uint8:
			token.Kind = KindUint8
			token.Value = uint8(value.Uint())
			return cont, nil

		case reflect.Uint16:
			token.Kind = KindUint16
			token.Value = uint16(value.Uint())
			return cont, nil

		case reflect.Uint32:
			token.Kind = KindUint32
			token.Value = uint32(value.Uint())
			return cont, nil

		case reflect.Uint64:
			token.Kind = KindUint64
			token.Value = uint64(value.Uint())
			return cont, nil

		case reflect.Uintptr:
			token.Kind = KindPointer
			token.Value = uintptr(value.Uint())
			return cont, nil

		case reflect.Float32:
			f := value.Float()
			if f != f {
				*token = NaN
				return cont, nil
			} else {
				token.Kind = KindFloat32
				token.Value = float32(f)
				return cont, nil
			}

		case reflect.Float64:
			f := value.Float()
			if f != f {
				*token = NaN
				return cont, nil
			} else {
				token.Kind = KindFloat64
				token.Value = f
				return cont, nil
			}

		case reflect.Array, reflect.Slice:
			if isBytes(value.Type()) {
				token.Kind = KindBytes
				token.Value = toBytes(value)
				return cont, nil
			} else {
				return MarshalArray(ctx, value, 0, cont), nil
			}

		case reflect.String:
			token.Kind = KindString
			token.Value = value.String()
			return cont, nil

		case reflect.Struct:
			return MarshalStruct(ctx, value, cont), nil

		case reflect.Map:
			return MarshalMap(ctx, value, cont), nil

		case reflect.Func:
			if ctx.IgnoreFuncs {
				*token = Nil
				return cont, nil
			}
			if value.Type().NumIn() != 0 {
				return nil, we.With(
					WithPath(ctx),
					e5.With(BadTupleType),
					e5.Info("bad tuple type: %v", value.Type()),
				)(
					MarshalError,
				)
			}
			var items []reflect.Value
			if !value.IsNil() {
				items = value.Call([]reflect.Value{})
			}
			return MarshalTuple(
				ctx,
				items,
				cont,
			), nil

		default:
			return cont, nil

		}
	}

	if value.IsValid() {
		if name, ok := registeredTypeToName.Load(value.Type()); ok {
			return func(token *Token) (Proc, error) {
				token.Kind = KindTypeName
				token.Value = name.(string)
				return marshal, nil
			}
		}
	}

	return marshal
}

var arrayEndToken = reflect.ValueOf(&Token{
	Kind: KindArrayEnd,
})

var arrayToken = Token{
	Kind: KindArray,
}

func MarshalArray(ctx Ctx, value reflect.Value, index int, cont Proc) Proc {
	var proc Proc
	proc = func(_ *Token) (Proc, error) {
		if index >= value.Len() {
			return ctx.Marshal(
				ctx,
				arrayEndToken,
				cont,
			), nil
		}
		v := value.Index(index)
		index++
		return ctx.Marshal(
			ctx.WithPath(index-1),
			v,
			proc,
		), nil
	}
	return func(token *Token) (Proc, error) {
		*token = arrayToken
		return proc, nil
	}
}

var objectEndToken = reflect.ValueOf(&Token{
	Kind: KindObjectEnd,
})

var objectToken = Token{
	Kind: KindObject,
}

func MarshalStruct(ctx Ctx, value reflect.Value, cont Proc) Proc {
	return func(token *Token) (Proc, error) {
		*token = objectToken
		return MarshalStructFields(ctx, value, cont), nil
	}
}

func MarshalStructFields(ctx Ctx, value reflect.Value, cont Proc) Proc {
	valueType := value.Type()
	numField := valueType.NumField()
	fieldIdx := 0
	var proc Proc
	proc = func(_ *Token) (Proc, error) {
		if fieldIdx == numField {
			return ctx.Marshal(
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
				return proc, nil
			}
			if field.Type.Kind() == reflect.Slice && fieldValue.Len() == 0 {
				fieldIdx++
				return proc, nil
			}
		}
		if field.PkgPath != "" {
			// unexported field
			fieldIdx++
			return proc, nil
		}

		if ctx.IgnoreFuncs && field.Type.Kind() == reflect.Func {
			fieldIdx++
			return proc, nil
		}

		fieldIdx++
		return ctx.Marshal(
			ctx.WithPath(field.Name),
			reflect.ValueOf(field.Name),
			func(token *Token) (Proc, error) {
				return ctx.Marshal(
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
		value.MapRange(),
		make([]*MapTuple, 0, value.Len()),
		cont,
	)
}

var mapEndToken = reflect.ValueOf(&Token{
	Kind: KindMapEnd,
})

var mapToken = Token{
	Kind: KindMap,
}

func MarshalMapIter(ctx Ctx, iter *reflect.MapIter, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func(_ *Token) (Proc, error) {
		if !iter.Next() {
			// done
			slices.SortFunc(tuples, func(a, b *MapTuple) int {
				return MustCompare(
					a.KeyTokens.Iter(),
					b.KeyTokens.Iter(),
				)
			})
			return MarshalMapTuples(ctx, tuples, cont), nil
		}
		var tokens Tokens
		keyMarshalProc := MarshalValue(Ctx{}, iter.Key(), nil)
		if err := Copy(
			// tokens are for sorting only, so do not call ctx.Marshal
			&keyMarshalProc,
			CollectTokens(&tokens),
		); err != nil {
			return nil, we.With(e5.With(MarshalError), WithPath(ctx))(err)
		}
		if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, we.With(WithPath(ctx), e5.With(BadMapKey))(MarshalError)
		}
		tuples = append(tuples, &MapTuple{
			KeyTokens: tokens,
			Key:       iter.Key(),
			Value:     iter.Value(),
		})
		return proc, nil
	}
	return func(token *Token) (Proc, error) {
		*token = mapToken
		return proc, nil
	}
}

func MarshalMapTuples(ctx Ctx, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func(_ *Token) (Proc, error) {
		if len(tuples) == 0 {
			return ctx.Marshal(
				ctx,
				mapEndToken,
				cont,
			), nil
		}
		tuple := tuples[0]
		tuples = tuples[1:]
		path := tuple.Key.Interface()
		return ctx.Marshal(
			ctx.WithPath(path),
			tuple.Key,
			// must wrap in closure to delay value marshaling
			func(token *Token) (Proc, error) {
				return ctx.Marshal(
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

var tupleToken = Token{
	Kind: KindTuple,
}

func MarshalTuple(ctx Ctx, items []reflect.Value, cont Proc) Proc {
	var proc Proc
	var n int
	proc = func(_ *Token) (Proc, error) {
		if len(items) == 0 {
			return ctx.Marshal(
				ctx,
				tupleEndToken,
				cont,
			), nil
		} else {
			v := items[0]
			items = items[1:]
			n++
			return ctx.Marshal(
				ctx.WithPath(n-1),
				v,
				proc,
			), nil
		}
	}
	return func(token *Token) (Proc, error) {
		*token = tupleToken
		return proc, nil
	}
}
