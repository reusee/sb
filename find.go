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
	anchor := root
	stack := []*Tree{
		root,
	}
	var token *Token
	var noHashNodes []*Tree

	for {
		token, err = stream.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		if token.Kind == KindPostHash {
			// set hash to anchor node
			if anchor.Token == nil {
				return nil, UnexpectedHashToken
			}
			if h, ok := token.Value.([]byte); ok {
				var node *Tree
				switch anchor.Kind {
				case KindArrayEnd,
					KindObjectEnd,
					KindMapEnd,
					KindTupleEnd:
					node = anchor.Paired
				default:
					node = anchor
				}
				node.Hash = h
				if bytes.Equal(h, hash) {
					subStream = node.Iter()
					return
				}
			}

		} else {
			if anchor != root && len(anchor.Hash) == 0 {
				// save for later rehash
				noHashNodes = append(noHashNodes, anchor)
			}
			node := &Tree{
				Token: token,
			}
			anchor = node
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

	// rehash
	if anchor != root && len(anchor.Hash) == 0 {
		noHashNodes = append(noHashNodes, anchor)
	}
	for _, node := range noHashNodes {
		var sum []byte
		sum, err = TreeHashSum(node, newState)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(sum, hash) {
			subStream = node.Iter()
			return
		}
	}

	if subStream == nil {
		err = fmt.Errorf("FindByHash %x: %w", hash, NotFound)
	}

	return
}
