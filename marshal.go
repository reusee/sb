package sb

import (
	"encoding"
	"math"
	"reflect"
	"sort"
)

type SBMarshaler interface {
	MarshalSB(ctx Ctx, cont Proc) Proc
}

var MarshalCtx = Ctx{
	Marshal: MarshalValue,
}

func Marshal(value any) *Proc {
	marshaler := MarshalValue(Ctx{
		Marshal: MarshalValue,
	}, reflect.ValueOf(value), nil)
	return &marshaler
}

func MarshalAny(ctx Ctx, value any, cont Proc) Proc {
	return ctx.Marshal(ctx, reflect.ValueOf(value), cont)
}

func MarshalValue(ctx Ctx, value reflect.Value, cont Proc) Proc {
	if ctx.Marshal == nil {
		ctx.Marshal = MarshalValue
	}
	return func() (*Token, Proc, error) {

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
				if math.IsNaN(v) {
					return NaN, cont, nil
				} else {
					return &Token{Kind: KindFloat64, Value: v}, cont, nil
				}

			case int8:
				return &Token{Kind: KindInt8, Value: v}, cont, nil

			case float32:
				if math.IsNaN(float64(v)) {
					return NaN, cont, nil
				} else {
					return &Token{Kind: KindFloat32, Value: v}, cont, nil
				}

			case int16:
				return &Token{Kind: KindInt16, Value: v}, cont, nil

			case SBMarshaler:
				return nil, v.MarshalSB(ctx, cont), nil

			case encoding.BinaryMarshaler:
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, nil, err
				}
				return nil, ctx.Marshal(ctx, reflect.ValueOf(string(bs)), cont), nil

			case encoding.TextMarshaler:
				bs, err := v.MarshalText()
				if err != nil {
					return nil, nil, err
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
			if math.IsNaN(value.Float()) {
				return NaN, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat32,
					Value: float32(value.Float()),
				}, cont, nil
			}

		case reflect.Float64:
			if math.IsNaN(value.Float()) {
				return NaN, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat64,
					Value: float64(value.Float()),
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
			ctx,
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
	var fields []reflect.StructField
	valueType := value.Type()
	numField := valueType.NumField()
	for i := 0; i < numField; i++ {
		field := valueType.Field(i)
		if ctx.SkipEmptyStructFields {
			fieldValue := value.Field(i)
			if fieldValue.IsZero() {
				continue
			}
			if field.Type.Kind() == reflect.Slice && fieldValue.Len() == 0 {
				continue
			}
		}
		if field.PkgPath == "" {
			// exported field
			fields = append(fields, field)
		}
	}
	if !ctx.ReserveStructFieldsOrder {
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})
	}
	return func() (*Token, Proc, error) {
		return objectToken, MarshalStructFields(ctx, value, fields, cont), nil
	}
}

func MarshalStructFields(ctx Ctx, value reflect.Value, fields []reflect.StructField, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(fields) == 0 {
			return nil, ctx.Marshal(
				ctx,
				objectEndToken,
				cont,
			), nil
		}
		field := fields[0]
		fields = fields[1:]
		return nil, ctx.Marshal(
			ctx,
			reflect.ValueOf(field.Name),
			func() (*Token, Proc, error) {
				return nil, ctx.Marshal(
					ctx,
					value.FieldByIndex(field.Index),
					proc,
				), nil
			},
		), nil
	}
	return proc
}

type MapTuple struct {
	KeyTokens Tokens
	Key       reflect.Value
	Value     reflect.Value
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
		if err := Copy(
			Marshal(iter.Key().Interface()),
			CollectTokens(&tokens),
		); err != nil {
			return nil, nil, err
		}
		if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, nil, MarshalError{BadMapKey}
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
		ctx.ReserveStructFieldsOrder = true
		return nil, ctx.Marshal(
			ctx,
			tuple.Key,
			func() (*Token, Proc, error) {
				ctx.ReserveStructFieldsOrder = false
				return nil, ctx.Marshal(
					ctx,
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
			return nil, ctx.Marshal(
				ctx,
				v,
				proc,
			), nil
		}
	}
	return func() (*Token, Proc, error) {
		return tupleToken, proc, nil
	}
}
