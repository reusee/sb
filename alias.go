package sb

import (
	"errors"
	"fmt"
)

type (
	any = interface{}
)

var (
	pt = fmt.Printf
	is = errors.Is
	as = errors.As
)
