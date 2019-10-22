package sb

type Tokenizer interface {
	Next() *Token
}

type Kind uint8

const (
	KindInvalid Kind = iota

	KindArrayEnd
	KindObjectEnd

	KindBool
	KindInt
	KindUint
	KindFloat
	KindString
	KindNil

	KindArray
	KindObject
)

type Token struct {
	Kind  Kind
	Value any
}
