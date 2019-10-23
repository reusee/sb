package sb

type Stream interface {
	Next() *Token
	Peek() *Token
}

type Kind uint8

const (
	KindInvalid Kind = iota

	KindArrayEnd
	KindObjectEnd

	KindNil
	KindBool
	KindString

	KindInt
	KindInt8
	KindInt16
	KindInt32
	KindInt64

	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64

	KindFloat32
	KindFloat64

	KindArray
	KindObject
)

type Token struct {
	Kind  Kind
	Value any
}
