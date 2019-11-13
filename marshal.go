package sb

import (
	"encoding"
	"math"
	"reflect"
	"sort"
	"sync"
)

type Marshaler struct {
	err  error
	proc Proc
}

type Proc func() (*Token, Proc)

var _ Stream = new(Marshaler)

func NewMarshaler(obj any) *Marshaler {
	m := new(Marshaler)
	m.proc = m.tokenize(
		reflect.ValueOf(obj),
		nil,
	)
	return m
}

type SBMarshaler interface {
	MarshalSB() []Token
}

func (t *Marshaler) generateTokens(tokens []Token, cont Proc) Proc {
	return func() (*Token, Proc) {
		if len(tokens) == 0 {
			return nil, cont
		}
		return &tokens[0], t.generateTokens(tokens[1:], cont)
	}
}

func (t *Marshaler) tokenize(value reflect.Value, cont Proc) Proc {
	return func() (*Token, Proc) {

		if value.IsValid() {
			i := value.Interface()
			if v, ok := i.(SBMarshaler); ok {
				return nil, t.generateTokens(v.MarshalSB(), cont)
			} else if v, ok := i.(encoding.BinaryMarshaler); ok {
				bs, err := v.MarshalBinary()
				if err != nil {
					t.err = err
					return nil, nil
				}
				return &Token{KindString, string(bs)}, cont
			} else if v, ok := i.(encoding.TextMarshaler); ok {
				bs, err := v.MarshalText()
				if err != nil {
					t.err = err
					return nil, nil
				}
				return &Token{KindString, string(bs)}, cont
			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			return &Token{
				Kind: KindNil,
			}, cont

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				return &Token{
					Kind: KindNil,
				}, cont
			} else {
				return nil, t.tokenize(value.Elem(), cont)
			}

		case reflect.Bool:
			return &Token{
				Kind:  KindBool,
				Value: bool(value.Bool()),
			}, cont

		case reflect.Int:
			return &Token{
				Kind:  KindInt,
				Value: int(value.Int()),
			}, cont

		case reflect.Int8:
			return &Token{
				Kind:  KindInt8,
				Value: int8(value.Int()),
			}, cont

		case reflect.Int16:
			return &Token{
				Kind:  KindInt16,
				Value: int16(value.Int()),
			}, cont

		case reflect.Int32:
			return &Token{
				Kind:  KindInt32,
				Value: int32(value.Int()),
			}, cont

		case reflect.Int64:
			return &Token{
				Kind:  KindInt64,
				Value: int64(value.Int()),
			}, cont

		case reflect.Uint:
			return &Token{
				Kind:  KindUint,
				Value: uint(value.Uint()),
			}, cont

		case reflect.Uint8:
			return &Token{
				Kind:  KindUint8,
				Value: uint8(value.Uint()),
			}, cont

		case reflect.Uint16:
			return &Token{
				Kind:  KindUint16,
				Value: uint16(value.Uint()),
			}, cont

		case reflect.Uint32:
			return &Token{
				Kind:  KindUint32,
				Value: uint32(value.Uint()),
			}, cont

		case reflect.Uint64:
			return &Token{
				Kind:  KindUint64,
				Value: uint64(value.Uint()),
			}, cont

		case reflect.Float32:
			if math.IsNaN(value.Float()) {
				return &Token{
					Kind: KindNaN,
				}, cont
			} else {
				return &Token{
					Kind:  KindFloat32,
					Value: float32(value.Float()),
				}, cont
			}

		case reflect.Float64:
			if math.IsNaN(value.Float()) {
				return &Token{
					Kind: KindNaN,
				}, cont
			} else {
				return &Token{
					Kind:  KindFloat64,
					Value: float64(value.Float()),
				}, cont
			}

		case reflect.Array, reflect.Slice:
			if isBytes(value.Type()) {
				return &Token{
					Kind:  KindBytes,
					Value: toBytes(value),
				}, cont
			} else {
				return &Token{
					Kind: KindArray,
				}, t.tokenizeArray(value, 0, cont)
			}

		case reflect.String:
			return &Token{
				Kind:  KindString,
				Value: value.String(),
			}, cont

		case reflect.Struct:
			return &Token{
				Kind: KindObject,
			}, t.tokenizeStruct(value, cont)

		case reflect.Map:
			return &Token{
				Kind: KindMap,
			}, t.tokenizeMap(value, cont)

		case reflect.Func:
			items := value.Call([]reflect.Value{})
			return &Token{
					Kind: KindTuple,
				}, t.tokenizeTuple(
					items,
					cont,
				)

		default:
			return nil, cont

		}
	}
}

func (t *Marshaler) tokenizeArray(value reflect.Value, index int, cont Proc) Proc {
	return func() (*Token, Proc) {
		if index >= value.Len() {
			return &Token{
				Kind: KindArrayEnd,
			}, cont
		}
		return nil, t.tokenize(
			value.Index(index),
			t.tokenizeArray(value, index+1, cont),
		)
	}
}

var structFields sync.Map

func (t *Marshaler) tokenizeStruct(value reflect.Value, cont Proc) Proc {
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
	return t.tokenizeStructField(value, fields, cont)
}

func (t *Marshaler) tokenizeStructField(value reflect.Value, fields []reflect.StructField, cont Proc) Proc {
	return func() (*Token, Proc) {
		if len(fields) == 0 {
			return &Token{
				Kind: KindObjectEnd,
			}, cont
		}
		field := fields[0]
		return &Token{
				Kind:  KindString,
				Value: field.Name,
			}, t.tokenize(
				value.FieldByIndex(field.Index),
				t.tokenizeStructField(value, fields[1:], cont),
			)
	}
}

type mapTuple struct {
	keyTokens Tokens
	value     reflect.Value
}

func (t *Marshaler) tokenizeMap(value reflect.Value, cont Proc) Proc {
	return t.tokenizeMapIter(
		value,
		value.MapRange(),
		make([]*mapTuple, 0, value.Len()),
		cont,
	)
}

func (t *Marshaler) tokenizeMapIter(
	value reflect.Value,
	iter *reflect.MapIter,
	tuples []*mapTuple,
	cont Proc,
) Proc {
	return func() (*Token, Proc) {
		if !iter.Next() {
			// done
			sort.Slice(tuples, func(i, j int) bool {
				return MustCompare(
					tuples[i].keyTokens.Iter(),
					tuples[j].keyTokens.Iter(),
				) < 0
			})
			return nil, t.tokenizeMapValue(tuples, cont)
		}
		tokens, err := TokensFromStream(NewMarshaler(iter.Key().Interface()))
		if err != nil {
			t.err = err
			return nil, nil
		} else if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			t.err = MarshalError{BadMapKey}
			return nil, nil
		}
		return nil, t.tokenizeMapIter(
			value,
			iter,
			append(tuples, &mapTuple{
				keyTokens: tokens,
				value:     iter.Value(),
			}),
			cont,
		)
	}
}

func (t *Marshaler) tokenizeMapValue(tuples []*mapTuple, cont Proc) Proc {
	return func() (*Token, Proc) {
		if len(tuples) == 0 {
			return &Token{
				Kind: KindMapEnd,
			}, cont
		}
		tuple := tuples[0]
		return nil, t.generateTokens(
			tuple.keyTokens,
			t.tokenize(
				tuple.value,
				t.tokenizeMapValue(tuples[1:], cont),
			),
		)
	}
}

func (t *Marshaler) tokenizeTuple(items []reflect.Value, cont Proc) Proc {
	return func() (*Token, Proc) {
		if len(items) == 0 {
			return &Token{
				Kind: KindTupleEnd,
			}, cont
		} else {
			return nil, t.tokenize(
				items[0],
				t.tokenizeTuple(
					items[1:],
					cont,
				),
			)
		}
	}
}

func (t *Marshaler) Next() (ret *Token, err error) {
check:
	if t.err != nil {
		return nil, t.err
	}
	if t.proc == nil {
		return nil, t.err
	}
	ret, t.proc = t.proc()
	if ret != nil {
		return ret, t.err
	}
	goto check
}
