package sb

import (
	"bytes"
	"fmt"
	"hash"
)

func FindByHash(
	stream Stream,
	hash []byte,
	newState func() hash.Hash,
) (
	subStream Stream,
	err error,
) {

	root := new(Tree)
	anchor := root
	stack := []*Tree{
		root,
	}
	var token *Token
	var noHashNodes []*Tree

	for {
		token, err = stream.Next()
		if err != nil { // NOCOVER
			return nil, err
		}
		if token == nil {
			break
		}

		if token.Kind == KindPostTag {
			// set to anchor node
			if anchor.Token == nil {
				return nil, UnexpectedHashToken
			}
			if tag, ok := token.Value.([]byte); ok {
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
				node.Tags.Add(tag)
				if bytes.HasPrefix(tag, []byte("hash:")) {
					if bytes.Equal(bytes.TrimPrefix(tag, []byte("hash:")), hash) {
						subStream = node.Iter()
						return
					}
				}
			}

		} else {
			if anchor != root {
				if _, ok := anchor.Tags.Get("hash"); !ok {
					// save for later rehash
					noHashNodes = append(noHashNodes, anchor)
				}
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
	if anchor != root {
		if _, ok := anchor.Tags.Get("hash"); !ok {
			noHashNodes = append(noHashNodes, anchor)
		}
	}
	for _, node := range noHashNodes {
		if err = node.FillHash(newState); err != nil { // NOCOVER
			return
		}
		h, ok := node.Tags.Get("hash")
		if !ok {
			panic("impossible")
		}
		if bytes.Equal(h, hash) {
			subStream = node.Iter()
			return
		}
	}

	if subStream == nil {
		err = fmt.Errorf("FindByHash %x: %w", hash, NotFound)
	}

	return
}
