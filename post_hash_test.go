package sb

import (
	"bytes"
	"hash/fnv"
	"os"
	"testing"
)

func TestHasher(t *testing.T) {
	for _, c := range marshalTestCases {
		marshaler := Marshal(c.value)
		hasher := NewPostHasher(marshaler, newMapHashState)
		tokens, err := TokensFromStream(hasher)
		if err != nil {
			t.Fatal(err)
		}
		hashToken := tokens[len(tokens)-1]
		if hashToken.Kind != KindPostTag {
			t.Fatal("not hash")
		}

		// compare hashed and not hashed
		if MustCompare(
			tokens.Iter(),
			Marshal(c.value),
		) != 0 {
			t.Fatal("not equal")
		}

		// hash tokens will be ignore
		hasher2 := NewPostHasher(tokens.Iter(), newMapHashState)
		tokens2, err := TokensFromStream(hasher2)
		if err != nil {
			t.Fatal(err)
		}
		hashToken2 := tokens2[len(tokens2)-1]
		if hashToken2.Kind != KindPostTag {
			t.Fatal("not hash")
		}
		if !bytes.Equal(hashToken.Value.([]byte), hashToken2.Value.([]byte)) {
			t.Fatal("hash not match")
		}

		// sum
		var sum1, sum2 []byte
		if err := Copy(
			Marshal(c.value),
			Hash(fnv.New128, &sum1, nil),
		); err != nil {
			t.Fatal(err)
		}
		if err := Copy(
			Marshal(c.value),
			Hash(fnv.New128a, &sum2, nil),
		); err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(sum1, sum2) {
			t.Fatal("should not equal")
		}
	}
}

func TestIntHash(t *testing.T) {
	tokens, err := TokensFromStream(
		NewPostHasher(
			Marshal(42),
			newMapHashState,
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 {
		dumpStream(tokens.Iter(), "->", os.Stdout)
		t.Fatalf("got %d\n", len(tokens))
	}
}
