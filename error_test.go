package sb

import (
	"errors"
	"io"
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
	if err.Error() != "DecodeError: string too long at 0" {
		t.Fatalf("got %s", err.Error())
	}

	err = NewMarshalError(DefaultCtx.WithPath("foo"), BadMapKey)
	if !errors.Is(err, BadMapKey) {
		t.Fatal()
	}
	if err.Error() != "MarshalError: bad map key at /foo" {
		t.Fatal()
	}

	var v map[int]map[int]string
	err = Copy(
		Marshal(map[int]map[int]int{
			42: {
				43: 44,
			},
		}),
		Unmarshal(&v),
	)
	if err.Error() != "UnmarshalError: expecting int at /42/43" {
		t.Fatal()
	}

	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Fatal()
			}
			if p != "bad data: string" {
				t.Fatalf("got %v", p)
			}
		}()
		NewDecodeError(0, io.EOF, "foo")
	}()

}
