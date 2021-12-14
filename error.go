package sb

import (
	"fmt"
	"reflect"
)

var (
	BadMapKey    = fmt.Errorf("bad map key")
	BadTokenKind = fmt.Errorf("bad token kind")

	UnexpectedEndToken = fmt.Errorf("unexpected end token")
)

// unmarshal

var UnmarshalError = fmt.Errorf("unmarshal error")

var (
	BadFieldName        = fmt.Errorf("bad field name")
	BadTargetType       = fmt.Errorf("bad target type")
	BadTupleType        = fmt.Errorf("bad tuple type")
	DuplicatedFieldName = fmt.Errorf("duplicated field name")
	TooManyElement      = fmt.Errorf("too many element")
	TooFewElement       = fmt.Errorf("too few element")
	UnknownFieldName    = fmt.Errorf("unknown field name")
)

type ErrUnmarshalTypeMismatch struct {
	TokenKind Kind
	Target    reflect.Kind
}

var _ error = ErrUnmarshalTypeMismatch{}

func (e ErrUnmarshalTypeMismatch) Error() string {
	return fmt.Sprintf("unmarshaling %s to %v", e.TokenKind, e.Target)
}

func TypeMismatch(kind Kind, target reflect.Kind) ErrUnmarshalTypeMismatch {
	return ErrUnmarshalTypeMismatch{
		TokenKind: kind,
		Target:    target,
	}
}

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
