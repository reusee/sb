package sb

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func TestSinkHash(t *testing.T) {
	type Case struct {
		value    any
		expected string
	}
	cases := []Case{
		{
			42,
			"0fcc339bcc03b2d67d97d0e2fa60bd41",
		},
		{
			[]int{1, 2, 3},
			"1686c4524aa5e66d9cf9b98296ea178c",
		},
	}

	for i, c := range cases {
		var sum []byte
		if err := Unmarshal(Marshal(c.value), Hasher(fnv.New128, &sum, nil)); err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sum) != c.expected {
			t.Fatalf("%d: %#v, got %x", i, c.value, sum)
		}

	}

}
