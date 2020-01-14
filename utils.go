package sb

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/maphash"
	"math/rand"
	"reflect"
)

type (
	any = interface{}
)

var (
	pt = fmt.Printf
	is = errors.Is

	anyType   = reflect.TypeOf((*any)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

func init() {
	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
}

var seed = maphash.MakeSeed()

func newMapHashState() hash.Hash {
	h := new(maphash.Hash)
	h.SetSeed(seed)
	return h
}
