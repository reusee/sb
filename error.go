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
	ExpectingString = fmt.Errorf("expecting string")
	ExpectingValue  = fmt.Errorf("expecting value")
)
