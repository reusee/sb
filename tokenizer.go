package sb

import (
	"fmt"
	"reflect"
)

func NewTokenizer(obj any) *Tokenizer {
	tokenizer := new(Tokenizer)
	tokenizer.proc = tokenizer.Tokenize(
		reflect.ValueOf(obj),
		tokenizer.End,
	)
	return tokenizer
}

type Tokenizer struct {
	proc   func()
	tokens []Token
}

func (t *Tokenizer) Tokenize(value reflect.Value, cont func()) func() {
	return func() {
		switch value.Kind() {

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
				Value: value.Bool(),
			})
			t.proc = cont

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindInt,
				Value: value.Int(),
			})
			t.proc = cont

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindUint,
				Value: value.Uint(),
			})
			t.proc = cont

		case reflect.Float32, reflect.Float64:
			t.tokens = append(t.tokens, Token{
				Kind:  KindFloat,
				Value: value.Float(),
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

func (t *Tokenizer) End() {
	t.proc = nil
}

func (t *Tokenizer) TokenizeArray(value reflect.Value, index int, cont func()) func() {
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

func (t *Tokenizer) TokenizeStruct(value reflect.Value, index int, cont func()) func() {
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
	tokenizer := NewTokenizer(obj)
	for tokenizer.proc != nil {
		tokenizer.proc()
	}
	return tokenizer.tokens
}

func (t *Tokenizer) Next() (ret *Token) {
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
