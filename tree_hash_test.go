package sb

import (
	"strings"
	"testing"
)

func TestBadTreeFillHash(t *testing.T) {

	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Fatal("should panic")
			}
			if p != "empty tree" {
				t.Fatal("not match")
			}
		}()
		new(Tree).FillHash(newMapHashState)
	}()

	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Fatal("should panic")
			}
			if !strings.HasPrefix(p.(error).Error(), "unexpected token") {
				t.Fatal("not match")
			}
		}()
		(&Tree{
			Token: &Token{
				Kind: 2,
			},
		}).FillHash(newMapHashState)
	}()

}
