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

var (
	ExpectingValue    = fmt.Errorf("expecting value")
	ExpectingString   = fmt.Errorf("expecting string")
	ExpectingBool     = fmt.Errorf("expecting bool")
	ExpectingInt      = fmt.Errorf("expecting int")
	ExpectingInt8     = fmt.Errorf("expecting int8")
	ExpectingInt16    = fmt.Errorf("expecting int16")
	ExpectingInt32    = fmt.Errorf("expecting int32")
	ExpectingInt64    = fmt.Errorf("expecting int64")
	ExpectingUint     = fmt.Errorf("expecting uint")
	ExpectingUint8    = fmt.Errorf("expecting uint8")
	ExpectingUint16   = fmt.Errorf("expecting uint16")
	ExpectingUint32   = fmt.Errorf("expecting uint32")
	ExpectingUint64   = fmt.Errorf("expecting uint64")
	ExpectingFloat32  = fmt.Errorf("expecting float32")
	ExpectingFloat64  = fmt.Errorf("expecting float64")
	ExpectingSequence = fmt.Errorf("expecting array / slice")
	ExpectingStruct   = fmt.Errorf("expecting struct")
)

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
	StringTooLong = fmt.Errorf("string too long")
)
