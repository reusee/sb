package sb

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRef(t *testing.T) {
	type foo struct {
		I int
		S string
	}
	tree := MustTreeFromStream(
		PostHash(
			Marshal(foo{
				I: 42,
				S: "42",
			}),
			newMapHashState,
		),
	)

	type ref struct {
		Hash  []byte
		Value any
	}
	var refs []ref

	refTree := tree.IterFunc(func(tree *Tree) (*Token, error) {
		if tree.Kind == KindString || tree.Kind == KindInt {
			refs = append(refs, ref{
				Hash:  tree.Hash,
				Value: tree.Value,
			})
			return &Token{
				Kind:  KindRef,
				Value: tree.Hash,
			}, nil
		}
		return nil, nil
	})

	n := 0
	deref := Deref(refTree, func(hash []byte) (Stream, error) {
		for _, ref := range refs {
			if bytes.Equal(ref.Hash, hash) {
				n++
				return Marshal(ref.Value), nil
			}
		}
		panic("ref not found")
	})
	if MustCompare(deref, tree.Iter()) != 0 {
		t.Fatal("not equal")
	}
	if n != 4 {
		t.Fatal("bad deref count")
	}
}

func TestBadDeref(t *testing.T) {
	tokens := Tokens{
		{
			Kind:  KindRef,
			Value: []byte{},
		},
	}
	s := Deref(tokens.Iter(), func(hash []byte) (Stream, error) {
		return nil, fmt.Errorf("foo")
	})
	_, err := TokensFromStream(s)
	if err.Error() != "foo" {
		t.Fatal()
	}
}
