package sb

import (
	"fmt"
	"hash"
)

var (
	MoreThanOneValue = fmt.Errorf("more than one value in stream")
)

type Tree struct {
	*Token
	Hash   []byte
	Subs   []*Tree
	Paired *Tree
}

type TreeOption interface {
	IsTreeOption()
}

type WithHash struct {
	NewHashState func() hash.Hash
}

func (_ WithHash) IsTreeOption() {}

func TreeFromStream(
	stream Stream,
	options ...TreeOption,
) (*Tree, error) {

	root := new(Tree)
	stack := []*Tree{
		root,
	}
	var hash []byte

	for _, option := range options {
		switch option := option.(type) {

		case WithHash:
			s := stream
			stream = Tee(s, HashFunc(
				option.NewHashState,
				nil,
				func(h []byte, _ *Token) error {
					if len(h) > 0 {
						hash = h
					}
					return nil
				},
				nil,
			))

		}
	}

	for {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			return nil, err
		}
		if token == nil {
			break
		}
		node := &Tree{
			Token: token,
			Hash:  hash,
		}
		parent := stack[len(stack)-1]
		parent.Subs = append(parent.Subs, node)
		switch token.Kind {
		case KindArray, KindObject, KindMap, KindTuple:
			stack = append(stack, node)
		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			if len(stack) == 1 {
				return nil, UnexpectedEndToken
			}
			node.Paired = parent
			stack = stack[:len(stack)-1]
		}
	}

	if len(root.Subs) > 1 {
		return nil, MoreThanOneValue
	}
	if len(root.Subs) == 1 {
		root = root.Subs[0]
	}
	root.Hash = hash

	return root, nil
}

func MustTreeFromStream(stream Stream, options ...TreeOption) *Tree {
	t, err := TreeFromStream(stream, options...)
	if err != nil { // NOCOVER
		panic(err)
	}
	return t
}
