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
	Paired *Tree
	Hash   []byte
	Subs   []*Tree
}

type TreeOption interface {
	IsTreeOption()
}

type WithHash struct {
	NewHashState func() hash.Hash
}

func (WithHash) IsTreeOption() {}

type TapTree struct {
	Func func(*Tree)
}

func (TapTree) IsTreeOption() {}

func TreeFromStream(
	stream Stream,
	options ...TreeOption,
) (*Tree, error) {

	root := new(Tree)
	stack := []*Tree{
		root,
	}
	var hash []byte
	var tap func(*Tree)

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
					} else {
						hash = nil
					}
					return nil
				},
				nil,
			))

		case TapTree:
			tap = option.Func

		}
	}

	for {
		var token Token
		err := stream.Next(&token)
		if err != nil { // NOCOVER
			return nil, err
		}
		if !token.Valid() {
			break
		}
		node := &Tree{
			Token: &token,
			Hash:  hash,
		}
		if tap != nil {
			tap(node)
		}
		parent := stack[len(stack)-1]
		if parent.Token != nil &&
			parent.Kind == KindTypeName &&
			len(parent.Subs) > 0 {
			// filled type name node
			stack = stack[:len(stack)-1]
			parent = stack[len(stack)-1]
		}
		parent.Subs = append(parent.Subs, node)
		switch token.Kind {
		case KindArray, KindObject, KindMap, KindTuple:
			stack = append(stack, node)
		case KindTypeName:
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
	if tap != nil {
		tap(root)
	}

	return root, nil
}

func MustTreeFromStream(stream Stream, options ...TreeOption) *Tree {
	t, err := TreeFromStream(stream, options...)
	if err != nil { // NOCOVER
		panic(err)
	}
	return t
}
