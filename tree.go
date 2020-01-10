package sb

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
	last := root
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
		if token.Kind == KindPostHash {
			// set hash to last node
			if last.Token == nil {
				return nil, UnexpectedHashToken
			}
			if hash, ok := token.Value.([]byte); ok {
				switch last.Kind {
				case KindArrayEnd,
					KindObjectEnd,
					KindMapEnd,
					KindTupleEnd:
					last.Hash = hash
					last.Paired.Hash = hash
				default:
					last.Hash = hash
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
	if len(root.Subs) == 1 {
		root = root.Subs[0]
	}
	return root, nil
}

func MustTreeFromStream(stream Stream) *Tree {
	t, err := TreeFromStream(stream)
	if err != nil {
		panic(err)
	}
	return t
}
