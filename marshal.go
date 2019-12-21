package sb

import (
	"encoding"
	"math"
	"reflect"
	"sort"
	"sync"
)

type MarshalProc func() (*Token, MarshalProc, error)

type SBMarshaler interface {
	MarshalSB(cont MarshalProc) MarshalProc
}

func NewMarshaler(value any) *MarshalProc {
	marshaler := MarshalValue(reflect.ValueOf(value), nil)
	return &marshaler
}

func MarshalTokens(tokens []Token, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if len(tokens) == 0 {
			return nil, cont, nil
		}
		return &tokens[0], MarshalTokens(tokens[1:], cont), nil
	}
}

func MarshalAny(value any, cont MarshalProc) MarshalProc {
	return MarshalValue(reflect.ValueOf(value), cont)
}

func MarshalValue(value reflect.Value, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {

		if value.IsValid() {
			i := value.Interface()
			if v, ok := i.(SBMarshaler); ok {
				return nil, v.MarshalSB(cont), nil
			} else if v, ok := i.(encoding.BinaryMarshaler); ok {
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, nil, err
				}
				return &Token{KindString, string(bs)}, cont, nil
			} else if v, ok := i.(encoding.TextMarshaler); ok {
				bs, err := v.MarshalText()
				if err != nil {
					return nil, nil, err
				}
				return &Token{KindString, string(bs)}, cont, nil
			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			return &Token{
				Kind: KindNil,
			}, cont, nil

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				return &Token{
					Kind: KindNil,
				}, cont, nil
			} else {
				return nil, MarshalValue(value.Elem(), cont), nil
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
				return &Token{
					Kind: KindNaN,
				}, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat32,
					Value: float32(value.Float()),
				}, cont, nil
			}

		case reflect.Float64:
			if math.IsNaN(value.Float()) {
				return &Token{
					Kind: KindNaN,
				}, cont, nil
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
				return &Token{
					Kind: KindArray,
				}, MarshalArray(value, 0, cont), nil
			}

		case reflect.String:
			return &Token{
				Kind:  KindString,
				Value: value.String(),
			}, cont, nil

		case reflect.Struct:
			return &Token{
				Kind: KindObject,
			}, MarshalStruct(value, cont), nil

		case reflect.Map:
			return &Token{
				Kind: KindMap,
			}, MarshalMap(value, cont), nil

		case reflect.Func:
			items := value.Call([]reflect.Value{})
			return &Token{
					Kind: KindTuple,
				}, MarshalTuple(
					items,
					cont,
				), nil

		default:
			return nil, cont, nil

		}
	}
}

func MarshalArray(value reflect.Value, index int, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if index >= value.Len() {
			return &Token{
				Kind: KindArrayEnd,
			}, cont, nil
		}
		return nil, MarshalValue(
			value.Index(index),
			MarshalArray(value, index+1, cont),
		), nil
	}
}

var structFields sync.Map

func MarshalStruct(value reflect.Value, cont MarshalProc) MarshalProc {
	var fields []reflect.StructField
	valueType := value.Type()
	if v, ok := structFields.Load(valueType); ok {
		fields = v.([]reflect.StructField)
	} else {
		numField := valueType.NumField()
		for i := 0; i < numField; i++ {
			fields = append(fields, valueType.Field(i))
		}
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})
		structFields.Store(valueType, fields)
	}
	return MarshalStructFields(value, fields, cont)
}

func MarshalStructFields(value reflect.Value, fields []reflect.StructField, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if len(fields) == 0 {
			return &Token{
				Kind: KindObjectEnd,
			}, cont, nil
		}
		field := fields[0]
		return &Token{
				Kind:  KindString,
				Value: field.Name,
			}, MarshalValue(
				value.FieldByIndex(field.Index),
				MarshalStructFields(value, fields[1:], cont),
			), nil
	}
}

type MapTuple struct {
	KeyTokens Tokens
	Value     reflect.Value
}

func MarshalMap(value reflect.Value, cont MarshalProc) MarshalProc {
	return MarshalMapIter(
		value,
		value.MapRange(),
		make([]*MapTuple, 0, value.Len()),
		cont,
	)
}

func MarshalMapIter(value reflect.Value, iter *reflect.MapIter, tuples []*MapTuple, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if !iter.Next() {
			// done
			sort.Slice(tuples, func(i, j int) bool {
				return MustCompare(
					tuples[i].KeyTokens.Iter(),
					tuples[j].KeyTokens.Iter(),
				) < 0
			})
			return nil, MarshalMapValue(tuples, cont), nil
		}
		marshaler := MarshalValue(iter.Key(), nil)
		tokens, err := TokensFromStream(&marshaler)
		if err != nil {
			return nil, nil, err
		} else if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, nil, MarshalError{BadMapKey}
		}
		return nil, MarshalMapIter(
			value,
			iter,
			append(tuples, &MapTuple{
				KeyTokens: tokens,
				Value:     iter.Value(),
			}),
			cont,
		), nil
	}
}

func MarshalMapValue(tuples []*MapTuple, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if len(tuples) == 0 {
			return &Token{
				Kind: KindMapEnd,
			}, cont, nil
		}
		tuple := tuples[0]
		return nil, MarshalTokens(
			tuple.KeyTokens,
			MarshalValue(
				tuple.Value,
				MarshalMapValue(tuples[1:], cont),
			),
		), nil
	}
}

func MarshalTuple(items []reflect.Value, cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		if len(items) == 0 {
			return &Token{
				Kind: KindTupleEnd,
			}, cont, nil
		} else {
			return nil, MarshalValue(
				items[0],
				MarshalTuple(
					items[1:],
					cont,
				),
			), nil
		}
	}
}

var _ Stream = (*MarshalProc)(nil)

func (p *MarshalProc) Next() (*Token, error) {
	for {
		if p == nil || *p == nil {
			return nil, nil
		}
		var ret *Token
		var err error
		ret, *p, err = (*p)()
		if ret != nil || err != nil {
			return ret, err
		}
	}
}
