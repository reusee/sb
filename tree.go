package sb

import (
	"fmt"
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

func TreeFromStream(
	stream Stream,
) (*Tree, error) {
	root := new(Tree)
	stack := []*Tree{
		root,
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
	return root, nil
}

func MustTreeFromStream(stream Stream) *Tree {
	t, err := TreeFromStream(stream)
	if err != nil { // NOCOVER
		panic(err)
	}
	return t
}
