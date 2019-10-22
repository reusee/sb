package sb

import (
	"fmt"
	"reflect"
)

type Value struct {
	proc   func()
	tokens []Token
}

var _ Tokenizer = new(Value)

func NewValue(obj any) *Value {
	tokenizer := new(Value)
	tokenizer.proc = tokenizer.Tokenize(
		reflect.ValueOf(obj),
		tokenizer.End,
	)
	return tokenizer
}

func (t *Value) Tokenize(value reflect.Value, cont func()) func() {
	return func() {
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
			t.tokens = append(t.tokens, Token{
				Kind:  KindFloat32,
				Value: float32(value.Float()),
			})
			t.proc = cont

		case reflect.Float64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindFloat64,
				Value: float64(value.Float()),
			})
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
			panic(fmt.Errorf("invalid value or type: %s", value.String()))

		}
	}
}

func (t *Value) End() {
	t.proc = nil
}

func (t *Value) TokenizeArray(value reflect.Value, index int, cont func()) func() {
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

func (t *Value) TokenizeStruct(value reflect.Value, index int, cont func()) func() {
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

func Tokens(obj any) []Token {
	tokenizer := NewValue(obj)
	for tokenizer.proc != nil {
		tokenizer.proc()
	}
	return tokenizer.tokens
}

func (t *Value) Next() (ret *Token) {
check:
	if len(t.tokens) > 0 {
		token := t.tokens[0]
		ret = &token
		t.tokens = append(t.tokens[:0], t.tokens[1:]...)
		return
	}
	if t.proc == nil {
		return nil
	}
	for len(t.tokens) == 0 && t.proc != nil {
		t.proc()
	}
	goto check
}
