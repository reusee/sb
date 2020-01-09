package sb

import (
	"bytes"
	"crypto/sha1"
	"hash/fnv"
	"testing"
)

func TestHasher(t *testing.T) {
	for _, c := range marshalTestCases {
		marshaler := NewMarshaler(c.value)
		hasher := NewHasher(marshaler, sha1.New)
		tokens, err := TokensFromStream(hasher)
		if err != nil {
			t.Fatal(err)
		}
		hashToken := tokens[len(tokens)-1]
		if hashToken.Kind != KindHash {
			t.Fatal("not hash")
		}

		// hash tokens will be ignore
		hasher2 := NewHasher(tokens.Iter(), sha1.New)
		tokens2, err := TokensFromStream(hasher2)
		if err != nil {
			t.Fatal(err)
		}
		hashToken2 := tokens2[len(tokens2)-1]
		if hashToken2.Kind != KindHash {
			t.Fatal("not hash")
		}
		if !bytes.Equal(hashToken.Value.([]byte), hashToken2.Value.([]byte)) {
			t.Fatal("hash not match")
		}

		// sum
		sum1, err := HashSum(NewMarshaler(c.value), fnv.New128)
		if err != nil {
			t.Fatal(err)
		}
		sum2, err := HashSum(NewMarshaler(c.value), fnv.New128a)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(sum1, sum2) {
			t.Fatal("should not equal")
		}
	}
}
