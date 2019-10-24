package sb

type Stream interface {
	Next() (*Token, error)
	Peek() (*Token, error)
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

	KindMax Kind = 0xFF
)

type Token struct {
	Kind  Kind
	Value any
}

func TokensFromStream(stream Stream) (tokens []Token, err error) {
	for {
		p, err := stream.Next()
		if err != nil {
			return nil, err
		}
		if p == nil {
			return tokens, err
		}
		tokens = append(tokens, *p)
	}
	return
}

func MustTokensFromStream(stream Stream) []Token {
	tokens, err := TokensFromStream(stream)
	if err != nil {
		panic(err)
	}
	return tokens
}
