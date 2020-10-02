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

	KindRef Kind = 251

	KindMax Kind = 0xFF
)
