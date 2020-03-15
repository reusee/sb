package sb

type Kind uint8

const (
	KindInvalid Kind = iota
	KindMin

	KindArrayEnd  Kind = 10
	KindObjectEnd Kind = 20
	KindMapEnd    Kind = 25
	KindTupleEnd  Kind = 27

	KindNil    Kind = 30
	KindBool   Kind = 40
	KindString Kind = 50
	KindBytes  Kind = 55

	KindInt   Kind = 60
	KindInt8  Kind = 70
	KindInt16 Kind = 80
	KindInt32 Kind = 90
	KindInt64 Kind = 100

	KindUint   Kind = 110
	KindUint8  Kind = 120
	KindUint16 Kind = 130
	KindUint32 Kind = 140
	KindUint64 Kind = 150

	KindFloat32 Kind = 160
	KindFloat64 Kind = 170
	KindNaN     Kind = 175

	KindArray  Kind = 180
	KindObject Kind = 190
	KindMap    Kind = 200
	KindTuple  Kind = 210

	KindPostTag Kind = 250
	KindRef     Kind = 251

	KindMax Kind = 0xFF
)

type Token struct {
	Kind  Kind
	Value any
}

type min struct{}

var _ SBMarshaler = min{}

func (_ min) MarshalSB(_ Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindMin,
		}, cont, nil
	}
}

var Min = min{}

type max struct{}

var _ SBMarshaler = max{}

func (_ max) MarshalSB(_ Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindMax,
		}, cont, nil
	}
}

var Max = max{}

func isEnd(k Kind) bool {
	switch k {
	case KindArrayEnd,
		KindObjectEnd,
		KindMapEnd,
		KindTupleEnd:
		return true
	}
	return false
}
