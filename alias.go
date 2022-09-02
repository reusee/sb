package sb

import (
	"errors"
	"fmt"

	"github.com/reusee/e5"
)

type (
	any = interface{}
)

var (
	pt = fmt.Printf
	is = errors.Is
	as = errors.As

	we = e5.Wrap
)
