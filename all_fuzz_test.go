//go:build gofuzzbeta
// +build gofuzzbeta

package sb

import (
	"testing"
)

func FuzzAll(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		t.Parallel()
		res := Fuzz(data)
		if res != 1 {
			t.Skip()
		}
	})
}
