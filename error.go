package sb

import "fmt"

type UnmarshalError struct {
	Prev error
}

func (u UnmarshalError) Unwrap() error {
	return u.Prev
}

func (u UnmarshalError) Error() string {
	return fmt.Sprintf("UnmarshalError: %s", u.Prev.Error())
}

type MarshalError struct {
	Prev error
}

func (u MarshalError) Unwrap() error {
	return u.Prev
}

func (u MarshalError) Error() string {
	return fmt.Sprintf("MarshalError: %s", u.Prev.Error())
}

type DecodeError struct {
	Prev error
}

func (d DecodeError) Unwrap() error {
	return d.Prev
}

func (d DecodeError) Error() string {
	return fmt.Sprintf("DecodeError: %s", d.Prev.Error())
}

var (
	ExpectingValue      = fmt.Errorf("expecting value")
	ExpectingString     = fmt.Errorf("expecting string")
	ExpectingBytes      = fmt.Errorf("expecting bytes")
	ExpectingBool       = fmt.Errorf("expecting bool")
	ExpectingInt        = fmt.Errorf("expecting int")
	ExpectingInt8       = fmt.Errorf("expecting int8")
	ExpectingInt16      = fmt.Errorf("expecting int16")
	ExpectingInt32      = fmt.Errorf("expecting int32")
	ExpectingInt64      = fmt.Errorf("expecting int64")
	ExpectingUint       = fmt.Errorf("expecting uint")
	ExpectingUint8      = fmt.Errorf("expecting uint8")
	ExpectingUint16     = fmt.Errorf("expecting uint16")
	ExpectingUint32     = fmt.Errorf("expecting uint32")
	ExpectingUint64     = fmt.Errorf("expecting uint64")
	ExpectingFloat32    = fmt.Errorf("expecting float32")
	ExpectingFloat64    = fmt.Errorf("expecting float64")
	ExpectingFloat      = fmt.Errorf("expecting float32 or float64")
	ExpectingSequence   = fmt.Errorf("expecting array / slice")
	ExpectingStruct     = fmt.Errorf("expecting struct")
	ExpectingMap        = fmt.Errorf("expecting map")
	ExpectingTuple      = fmt.Errorf("expecting tuple")
	BadFieldName        = fmt.Errorf("bad field name")
	DuplicatedFieldName = fmt.Errorf("duplicated field name")
	TooManyElement      = fmt.Errorf("too many element")
	BadMapKey           = fmt.Errorf("bad map key")
	StringTooLong       = fmt.Errorf("string too long")
	BytesTooLong        = fmt.Errorf("bytes too long")
	BadTokenKind        = fmt.Errorf("bad token kind")
	BadTupleType        = fmt.Errorf("bad tuple type")
	BadTargetType       = fmt.Errorf("bad target type")
	UnexpectedEndToken  = fmt.Errorf("unexpected end token")
)
