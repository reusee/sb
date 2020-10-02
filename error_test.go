package sb

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	var err error
	err = NewUnmarshalError(DefaultCtx, ExpectingInt)
	if !errors.Is(err, ExpectingInt) {
		t.Fatal()
	}
	if err.Error() != "UnmarshalError: expecting int" {
		t.Fatalf("got %s", err.Error())
	}
	err = NewDecodeError(0, StringTooLong)
	if !errors.Is(err, StringTooLong) {
		t.Fatal()
	}
	if err.Error() != "DecodeError: string too long" {
		t.Fatal()
	}
	err = NewMarshalError(DefaultCtx, BadMapKey)
	if !errors.Is(err, BadMapKey) {
		t.Fatal()
	}
	if err.Error() != "MarshalError: bad map key" {
		t.Fatal()
	}
}
