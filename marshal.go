package sb

import (
	"encoding"
	"math"
	"reflect"
	"sort"
	"sync"
)

type Marshaler struct {
	tokens []Token
	err    error
	proc   Proc
}

type Proc func() Proc

var _ Stream = new(Marshaler)

func NewMarshaler(obj any) *Marshaler {
	m := new(Marshaler)
	m.proc = m.Tokenize(
		reflect.ValueOf(obj),
		nil,
	)
	return m
}

type Tokenizer interface {
	TokenizeSB() []Token
}

func (t *Marshaler) Tokenize(value reflect.Value, cont Proc) Proc {
	return func() Proc {

		if value.IsValid() {
			i := value.Interface()
			if v, ok := i.(Tokenizer); ok {
				t.tokens = append(t.tokens, v.TokenizeSB()...)
				return cont
			} else if v, ok := i.(encoding.BinaryMarshaler); ok {
				bs, err := v.MarshalBinary()
				if err != nil {
					t.err = err
					return nil
				}
				t.tokens = append(t.tokens, Token{KindString, string(bs)})
				return cont
			} else if v, ok := i.(encoding.TextMarshaler); ok {
				bs, err := v.MarshalText()
				if err != nil {
					t.err = err
					return nil
				}
				t.tokens = append(t.tokens, Token{KindString, string(bs)})
				return cont
			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			t.tokens = append(t.tokens, Token{
				Kind: KindNil,
			})
			return cont

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				t.tokens = append(t.tokens, Token{
					Kind: KindNil,
				})
				return cont
			} else {
				return t.Tokenize(value.Elem(), cont)
			}

		case reflect.Bool:
			t.tokens = append(t.tokens, Token{
				Kind:  KindBool,
				Value: bool(value.Bool()),
			})
			return cont

		case reflect.Int:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt,
				Value: int(value.Int()),
			})
			return cont

		case reflect.Int8:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt8,
				Value: int8(value.Int()),
			})
			return cont

		case reflect.Int16:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt16,
				Value: int16(value.Int()),
			})
			return cont

		case reflect.Int32:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt32,
				Value: int32(value.Int()),
			})
			return cont

		case reflect.Int64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt64,
				Value: int64(value.Int()),
			})
			return cont

		case reflect.Uint:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint,
				Value: uint(value.Uint()),
			})
			return cont

		case reflect.Uint8:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint8,
				Value: uint8(value.Uint()),
			})
			return cont

		case reflect.Uint16:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint16,
				Value: uint16(value.Uint()),
			})
			return cont

		case reflect.Uint32:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint32,
				Value: uint32(value.Uint()),
			})
			return cont

		case reflect.Uint64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint64,
				Value: uint64(value.Uint()),
			})
			return cont

		case reflect.Float32:
			if math.IsNaN(value.Float()) {
				t.tokens = append(t.tokens, Token{
					Kind: KindNaN,
				})
			} else {
				t.tokens = append(t.tokens, Token{
					Kind:  KindFloat32,
					Value: float32(value.Float()),
				})
			}
			return cont

		case reflect.Float64:
			if math.IsNaN(value.Float()) {
				t.tokens = append(t.tokens, Token{
					Kind: KindNaN,
				})
			} else {
				t.tokens = append(t.tokens, Token{
					Kind:  KindFloat64,
					Value: float64(value.Float()),
				})
			}
			return cont

		case reflect.Array, reflect.Slice:
			if isBytes(value.Type()) {
				t.tokens = append(t.tokens, Token{
					Kind:  KindBytes,
					Value: toBytes(value),
				})
				return cont
			} else {
				t.tokens = append(t.tokens, Token{
					Kind: KindArray,
				})
				return t.TokenizeArray(value, 0, cont)
			}

		case reflect.String:
			t.tokens = append(t.tokens, Token{
				Kind:  KindString,
				Value: value.String(),
			})
			return cont

		case reflect.Struct:
			t.tokens = append(t.tokens, Token{
				Kind: KindObject,
			})
			return t.TokenizeStruct(value, cont)

		case reflect.Map:
			t.tokens = append(t.tokens, Token{
				Kind: KindMap,
			})
			return t.TokenizeMap(value, cont)

		default:
			return cont

		}
	}
}

func (t *Marshaler) TokenizeArray(value reflect.Value, index int, cont Proc) Proc {
	return func() Proc {
		if index >= value.Len() {
			t.tokens = append(t.tokens, Token{
				Kind: KindArrayEnd,
			})
			return cont
		}
		return t.Tokenize(
			value.Index(index),
			t.TokenizeArray(value, index+1, cont),
		)
	}
}

var structFields sync.Map

func (t *Marshaler) TokenizeStruct(value reflect.Value, cont Proc) Proc {
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
	return t.TokenizeStructField(value, fields, cont)
}

func (t *Marshaler) TokenizeStructField(value reflect.Value, fields []reflect.StructField, cont Proc) Proc {
	return func() Proc {
		if len(fields) == 0 {
			t.tokens = append(t.tokens, Token{
				Kind: KindObjectEnd,
			})
			return cont
		}
		field := fields[0]
		t.tokens = append(t.tokens, Token{
			Kind:  KindString,
			Value: field.Name,
		})
		return t.Tokenize(
			value.FieldByIndex(field.Index),
			t.TokenizeStructField(value, fields[1:], cont),
		)
	}
}

type mapTuple struct {
	keyTokens Tokens
	value     reflect.Value
}

func (t *Marshaler) TokenizeMap(value reflect.Value, cont Proc) Proc {
	return t.TokenizeMapIter(value, value.MapRange(), []mapTuple{}, cont)
}

func (t *Marshaler) TokenizeMapIter(
	value reflect.Value,
	iter *reflect.MapIter,
	tuples []mapTuple,
	cont Proc,
) Proc {
	return func() Proc {
		if !iter.Next() {
			// done
			sort.Slice(tuples, func(i, j int) bool {
				return MustCompare(
					tuples[i].keyTokens.Iter(),
					tuples[j].keyTokens.Iter(),
				) < 0
			})
			return t.TokenizeMapValue(tuples, cont)
		}
		tokens, err := TokensFromStream(NewMarshaler(iter.Key().Interface()))
		if err != nil {
			t.err = err
			return nil
		} else if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			t.err = MarshalError{BadMapKey}
			return nil
		}
		return t.TokenizeMapIter(
			value,
			iter,
			append(tuples, mapTuple{
				keyTokens: tokens,
				value:     iter.Value(),
			}),
			cont,
		)
	}
}

func (t *Marshaler) TokenizeMapValue(tuples []mapTuple, cont Proc) Proc {
	return func() Proc {
		if len(tuples) == 0 {
			t.tokens = append(t.tokens, Token{
				Kind: KindMapEnd,
			})
			return cont
		}
		tuple := tuples[0]
		t.tokens = append(t.tokens, tuple.keyTokens...)
		return t.Tokenize(
			tuple.value,
			t.TokenizeMapValue(tuples[1:], cont),
		)
	}
}

func (t *Marshaler) Next() (ret *Token, err error) {
	ret, err = t.Peek()
	if ret != nil {
		t.tokens = append(t.tokens[:0], t.tokens[1:]...)
	}
	return
}

func (t *Marshaler) Peek() (ret *Token, err error) {
check:
	if t.err != nil {
		return nil, t.err
	}
	if len(t.tokens) > 0 {
		return &t.tokens[0], nil
	}
	if t.proc == nil {
		return nil, t.err
	}
	for len(t.tokens) == 0 && t.proc != nil {
		t.proc = t.proc()
	}
	goto check
}
