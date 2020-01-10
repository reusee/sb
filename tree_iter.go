package sb

func (t *Tree) Iter() *Proc {
	proc := IterTree(t, nil)
	return &proc
}

func IterTree(
	tree *Tree,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		return tree.Token, IterSubTrees(
			tree.Subs, 0,
			func() (*Token, Proc, error) {
				if len(tree.Hash) > 0 {
					return &Token{
						Kind:  KindPostHash,
						Value: tree.Hash,
					}, cont, nil
				}
				return nil, cont, nil
			},
		), nil
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
