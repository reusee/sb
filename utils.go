package sb

import (
	crand "crypto/rand"
	"encoding/binary"
	"hash"
	"hash/maphash"
	"math/rand"
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
