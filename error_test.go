package sb

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	var err error
	err = UnmarshalError{ExpectingInt}
	if !errors.Is(err, ExpectingInt) {
		t.Fatal()
	}
	if err.Error() != "UnmarshalError: expecting int" {
		t.Fatal()
	}
	err = DecodeError{StringTooLong}
	if !errors.Is(err, StringTooLong) {
		t.Fatal()
	}
	if err.Error() != "DecodeError: string too long" {
		t.Fatal()
	}
	err = MarshalError{BadMapKey}
	if !errors.Is(err, BadMapKey) {
		t.Fatal()
	}
	if err.Error() != "MarshalError: bad map key" {
		t.Fatal()
	}
}
