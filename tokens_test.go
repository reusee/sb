package sb

import (
	"bytes"
	"testing"
)

func TestTokens(t *testing.T) {
	// bad stream
	func() {
		defer func() {
			if p := recover(); p == nil {
				t.Fatal()
			}
		}()
		MustTokensFromStream(Decode(bytes.NewReader([]byte{
			byte(KindString), // incomplete
		})))
	}()

	for _, c := range marshalTestCases {
		MustTokensFromStream(Marshal(c.value))
	}

}
