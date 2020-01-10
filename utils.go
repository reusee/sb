package sb

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
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

func init() {
	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
}
