package sb

type Stream interface {
	Next() (*Token, error)
}

type Kind uint8

const (
	KindInvalid Kind = iota
	KindMin

	KindArrayEnd  = 10
	KindObjectEnd = 20
	KindMapEnd    = 25
	KindTupleEnd  = 27

	KindNil    = 30
	KindBool   = 40
	KindString = 50
	KindBytes  = 55

	KindInt   = 60
	KindInt8  = 70
	KindInt16 = 80
	KindInt32 = 90
	KindInt64 = 100

	KindUint   = 110
	KindUint8  = 120
	KindUint16 = 130
	KindUint32 = 140
	KindUint64 = 150

	KindFloat32 = 160
	KindFloat64 = 170
	KindNaN     = 175

	KindArray  = 180
	KindObject = 190
	KindMap    = 200
	KindTuple  = 210

	KindHash = 250

	KindMax Kind = 0xFF
)

type Token struct {
	Kind  Kind
	Value any
}

type min struct{}

var _ SBMarshaler = min{}

func (_ min) MarshalSB(cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		return &Token{
			Kind: KindMin,
		}, cont, nil
	}
}

var Min = min{}

type max struct{}

var _ SBMarshaler = max{}

func (_ max) MarshalSB(cont MarshalProc) MarshalProc {
	return func() (*Token, MarshalProc, error) {
		return &Token{
			Kind: KindMax,
		}, cont, nil
	}
}

var Max = max{}
