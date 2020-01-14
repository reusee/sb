package sb

import (
	"bytes"
	"fmt"
	"hash"
)

type HashFinder interface {
	FindByHash(
		hash []byte,
		newState func() hash.Hash,
	) (
		subStream Stream,
		err error,
	)
}

func FindByHash(
	stream Stream,
	hash []byte,
	newState func() hash.Hash,
) (
	subStream Stream,
	err error,
) {

	if finder, ok := stream.(HashFinder); ok {
		return finder.FindByHash(hash, newState)
	}

	// build tree
	root := new(Tree)
	last := root
	stack := []*Tree{
		root,
	}
	var token *Token
	for {
		token, err = stream.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		if token.Kind == KindPostHash {
			// set hash to last node
			if last.Token == nil {
				return nil, UnexpectedHashToken
			}
			if h, ok := token.Value.([]byte); ok {
				var node *Tree
				switch last.Kind {
				case KindArrayEnd,
					KindObjectEnd,
					KindMapEnd,
					KindTupleEnd:
					node = last.Paired
				default:
					node = last
				}
				node.Hash = h
				if bytes.Equal(h, hash) {
					subStream = node.Iter()
					return
				}
			}

		} else {
			node := &Tree{
				Token: token,
			}
			last = node
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

	}

	if len(root.Subs) > 1 {
		return nil, MoreThanOneValue
	}

	if subStream == nil {
		err = fmt.Errorf("FindByHash %x: %w", hash, NotFound)
	}

	return
}
