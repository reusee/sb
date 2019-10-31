package sb

import (
	"encoding"
	"math"
	"reflect"
)

type Marshaler struct {
	proc   func()
	tokens []Token
	err    error
}

var _ Stream = new(Marshaler)

func NewMarshaler(obj any) *Marshaler {
	m := new(Marshaler)
	m.proc = m.Tokenize(
		reflect.ValueOf(obj),
		m.End,
	)
	return m
}

type Tokenizer interface {
	TokenizeSB() []Token
}

var (
	tokenizerType       = reflect.TypeOf((*Tokenizer)(nil)).Elem()
	binaryMarshalerType = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	textMarshalerType   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func (t *Marshaler) Tokenize(value reflect.Value, cont func()) func() {
	return func() {

		if value.IsValid() {
			if value.Type().Implements(tokenizerType) {
				t.tokens = append(t.tokens, value.Interface().(Tokenizer).TokenizeSB()...)
				t.proc = cont
				return
			} else if value.Type().Implements(binaryMarshalerType) {
				bs, err := value.Interface().(encoding.BinaryMarshaler).MarshalBinary()
				if err != nil {
					t.err = err
					t.proc = nil
					return
				}
				t.tokens = append(t.tokens, Token{KindString, string(bs)})
				t.proc = cont
				return
			} else if value.Type().Implements(textMarshalerType) {
				bs, err := value.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					t.err = err
					t.proc = nil
					return
				}
				t.tokens = append(t.tokens, Token{KindString, string(bs)})
				t.proc = cont
				return
			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			t.tokens = append(t.tokens, Token{
				Kind: KindNil,
			})
			t.proc = cont

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				t.tokens = append(t.tokens, Token{
					Kind: KindNil,
				})
				t.proc = cont
			} else {
				t.proc = t.Tokenize(value.Elem(), cont)
			}

		case reflect.Bool:
			t.tokens = append(t.tokens, Token{
				Kind:  KindBool,
				Value: bool(value.Bool()),
			})
			t.proc = cont

		case reflect.Int:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt,
				Value: int(value.Int()),
			})
			t.proc = cont

		case reflect.Int8:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt8,
				Value: int8(value.Int()),
			})
			t.proc = cont

		case reflect.Int16:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt16,
				Value: int16(value.Int()),
			})
			t.proc = cont

		case reflect.Int32:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt32,
				Value: int32(value.Int()),
			})
			t.proc = cont

		case reflect.Int64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt64,
				Value: int64(value.Int()),
			})
			t.proc = cont

		case reflect.Uint:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint,
				Value: uint(value.Uint()),
			})
			t.proc = cont

		case reflect.Uint8:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint8,
				Value: uint8(value.Uint()),
			})
			t.proc = cont

		case reflect.Uint16:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint16,
				Value: uint16(value.Uint()),
			})
			t.proc = cont

		case reflect.Uint32:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint32,
				Value: uint32(value.Uint()),
			})
			t.proc = cont

		case reflect.Uint64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint64,
				Value: uint64(value.Uint()),
			})
			t.proc = cont

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
			t.proc = cont

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
			t.proc = cont

		case reflect.Array, reflect.Slice:
			t.tokens = append(t.tokens, Token{
				Kind: KindArray,
			})
			t.proc = t.TokenizeArray(value, 0, cont)

		case reflect.String:
			t.tokens = append(t.tokens, Token{
				Kind:  KindString,
				Value: value.String(),
			})
			t.proc = cont

		case reflect.Struct:
			t.tokens = append(t.tokens, Token{
				Kind: KindObject,
			})
			t.proc = t.TokenizeStruct(value, 0, cont)

		default:
			t.proc = cont

		}
	}
}

func (t *Marshaler) End() {
	t.proc = nil
}

func (t *Marshaler) TokenizeArray(value reflect.Value, index int, cont func()) func() {
	return func() {
		if index >= value.Len() {
			t.tokens = append(t.tokens, Token{
				Kind: KindArrayEnd,
			})
			t.proc = cont
			return
		}
		t.proc = t.Tokenize(
			value.Index(index),
			t.TokenizeArray(value, index+1, cont),
		)
	}
}

func (t *Marshaler) TokenizeStruct(value reflect.Value, index int, cont func()) func() {
	return func() {
		if index >= value.NumField() {
			t.tokens = append(t.tokens, Token{
				Kind: KindObjectEnd,
			})
			t.proc = cont
			return
		}
		t.tokens = append(t.tokens, Token{
			Kind:  KindString,
			Value: value.Type().Field(index).Name,
		})
		t.proc = t.Tokenize(
			value.Field(index),
			t.TokenizeStruct(value, index+1, cont),
		)
	}
}

func Tokens(obj any) ([]Token, error) {
	m := NewMarshaler(obj)
	for m.proc != nil {
		m.proc()
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.tokens, nil
}

func (t *Marshaler) Next() (ret *Token, err error) {
check:
	if t.err != nil {
		return nil, err
	}
	if len(t.tokens) > 0 {
		token := t.tokens[0]
		ret = &token
		t.tokens = append(t.tokens[:0], t.tokens[1:]...)
		return
	}
	if t.proc == nil {
		return nil, nil
	}
	for len(t.tokens) == 0 && t.proc != nil {
		t.proc()
	}
	goto check
}

func (t *Marshaler) Peek() (ret *Token, err error) {
check:
	if t.err != nil {
		return nil, err
	}
	if len(t.tokens) > 0 {
		token := t.tokens[0]
		return &token, nil
	}
	if t.proc == nil {
		return nil, nil
	}
	for len(t.tokens) == 0 && t.proc != nil {
		t.proc()
	}
	goto check
}
