package sb

func (t *Tree) Iter() *Proc {
	proc := IterTree(t, nil)
	return &proc
}

func IterTree(
	tree *Tree,
	cont Proc,
) Proc {
	return func() (token *Token, next Proc, err error) {
		token = tree.Token
		next = IterSubTrees(tree.Subs, 0, cont)
		return
	}
}

func IterSubTrees(
	subs []*Tree,
	index int,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		if len(subs) == 0 {
			return nil, cont, nil
		}
		if index >= len(subs) {
			return nil, cont, nil
		}
		return nil, IterTree(
			subs[index],
			IterSubTrees(
				subs,
				index+1,
				cont,
			),
		), nil
	}
}
