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
				if !isEnd(tree.Kind) && len(tree.Hash) > 0 {
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

func (t *Tree) IterFunc(
	fn func(*Tree) (*Token, error),
) *Proc {
	proc := IterTreeFunc(t, fn, nil)
	return &proc
}

func IterTreeFunc(
	tree *Tree,
	fn func(*Tree) (*Token, error),
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		token, err := fn(tree)
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token != nil {
			return token, cont, nil
		}
		return tree.Token, IterSubTreesFunc(
			tree.Subs, 0, fn,
			func() (*Token, Proc, error) {
				if !isEnd(tree.Kind) && len(tree.Hash) > 0 {
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

func IterSubTreesFunc(
	subs []*Tree,
	index int,
	fn func(*Tree) (*Token, error),
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		if len(subs) == 0 {
			return nil, cont, nil
		}
		if index >= len(subs) {
			return nil, cont, nil
		}
		return nil, IterTreeFunc(
			subs[index],
			fn,
			IterSubTreesFunc(
				subs,
				index+1,
				fn,
				cont,
			),
		), nil
	}
}
