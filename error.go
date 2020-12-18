package sb

import (
	"fmt"
)

var (
	BadMapKey    = fmt.Errorf("bad map key")
	BadTokenKind = fmt.Errorf("bad token kind")

	UnexpectedEndToken = fmt.Errorf("unexpected end token")

	ExpectingArrayEnd  = fmt.Errorf("expecting array end")
	ExpectingBool      = fmt.Errorf("expecting bool")
	ExpectingBytes     = fmt.Errorf("expecting bytes")
	ExpectingFloat     = fmt.Errorf("expecting float32 or float64")
	ExpectingFloat32   = fmt.Errorf("expecting float32")
	ExpectingFloat64   = fmt.Errorf("expecting float64")
	ExpectingInt       = fmt.Errorf("expecting int")
	ExpectingInt16     = fmt.Errorf("expecting int16")
	ExpectingInt32     = fmt.Errorf("expecting int32")
	ExpectingInt64     = fmt.Errorf("expecting int64")
	ExpectingInt8      = fmt.Errorf("expecting int8")
	ExpectingMap       = fmt.Errorf("expecting map")
	ExpectingMapEnd    = fmt.Errorf("expecting map end")
	ExpectingObjectEnd = fmt.Errorf("expecting object end")
	ExpectingSequence  = fmt.Errorf("expecting array / slice")
	ExpectingString    = fmt.Errorf("expecting string")
	ExpectingStruct    = fmt.Errorf("expecting struct")
	ExpectingTuple     = fmt.Errorf("expecting tuple")
	ExpectingTupleEnd  = fmt.Errorf("expecting tuple end")
	ExpectingUint      = fmt.Errorf("expecting uint")
	ExpectingUint16    = fmt.Errorf("expecting uint16")
	ExpectingUint32    = fmt.Errorf("expecting uint32")
	ExpectingUint64    = fmt.Errorf("expecting uint64")
	ExpectingUint8     = fmt.Errorf("expecting uint8")
	ExpectingValue     = fmt.Errorf("expecting value")

	kindToExpectingErr = map[Kind]error{
		KindString:    ExpectingString,
		KindBytes:     ExpectingBytes,
		KindBool:      ExpectingBool,
		KindInt:       ExpectingInt,
		KindInt8:      ExpectingInt8,
		KindInt16:     ExpectingInt16,
		KindInt32:     ExpectingInt32,
		KindInt64:     ExpectingInt64,
		KindUint:      ExpectingUint,
		KindUint8:     ExpectingUint8,
		KindUint16:    ExpectingUint16,
		KindUint32:    ExpectingUint32,
		KindUint64:    ExpectingUint64,
		KindFloat32:   ExpectingFloat32,
		KindFloat64:   ExpectingFloat64,
		KindArray:     ExpectingSequence,
		KindObject:    ExpectingStruct,
		KindMap:       ExpectingMap,
		KindTuple:     ExpectingTuple,
		KindArrayEnd:  ExpectingArrayEnd,
		KindObjectEnd: ExpectingObjectEnd,
		KindMapEnd:    ExpectingMapEnd,
		KindTupleEnd:  ExpectingTupleEnd,
	}
)

// unmarshal

var UnmarshalError = fmt.Errorf("unmarshal error")

var (
	BadFieldName        = fmt.Errorf("bad field name")
	BadTargetType       = fmt.Errorf("bad target type")
	BadTupleType        = fmt.Errorf("bad tuple type")
	DuplicatedFieldName = fmt.Errorf("duplicated field name")
	TooManyElement      = fmt.Errorf("too many element")
	UnknownFieldName    = fmt.Errorf("unknown field name")
)

// marshal

var MarshalError = fmt.Errorf("marshal error")

var (
	CyclicPointer = fmt.Errorf("cyclic pointer")
)

// decode

var DecodeError = fmt.Errorf("decode error")

var (
	StringTooLong   = fmt.Errorf("string too long")
	BytesTooLong    = fmt.Errorf("bytes too long")
	BadStringLength = fmt.Errorf("bad string length")
)
