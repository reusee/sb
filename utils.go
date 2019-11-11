package sb

import (
	"fmt"
	"reflect"
)

type (
	any = interface{}
)

var (
	pt = fmt.Printf

	anyType   = reflect.TypeOf((*any)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)
