package sb

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestTreeFillHash(t *testing.T) {
	tree := MustTreeFromStream(Marshal(42))
	if err := tree.FillHash(sha256.New); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", tree.Hash) != "151a3a0b4c88483512fc484d0badfedf80013ebb18df498bbee89ac5b69d7222" {
		t.Fatalf("got %x", tree.Hash)
	}
	if err := tree.FillHash(sha256.New); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", tree.Hash) != "151a3a0b4c88483512fc484d0badfedf80013ebb18df498bbee89ac5b69d7222" {
		t.Fatalf("got %x", tree.Hash)
	}

	for _, c := range marshalTestCases {
		tree := MustTreeFromStream(MarshalCtx(c.ctx, c.value))
		if err := tree.FillHash(sha256.New); err != nil {
			t.Fatal(err)
		}
		var hash []byte
		if err := Copy(
			MarshalCtx(c.ctx, c.value),
			Hash(sha256.New, &hash, nil),
		); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(hash, tree.Hash) {
			t.Fatal()
		}
	}
}

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
