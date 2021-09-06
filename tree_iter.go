package sb

func (t *Tree) Iter() Proc {
	return IterTree(t, nil)
}

func IterTree(
	tree *Tree,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		return tree.Token, IterSubTrees(
			tree.Subs, 0,
			func() (*Token, Proc, error) {
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
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(subs) == 0 {
			return nil, cont, nil
		}
		if index >= len(subs) {
			return nil, cont, nil
		}
		sub := subs[index]
		index++
		return nil, IterTree(
			sub,
			proc,
		), nil
	}
	return proc
}

func (t *Tree) IterFunc(
	fn func(*Tree) (*Token, error),
) Proc {
	return IterTreeFunc(t, fn, nil)
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
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(subs) == 0 {
			return nil, cont, nil
		}
		if index >= len(subs) {
			return nil, cont, nil
		}
		sub := subs[index]
		index++
		return nil, IterTreeFunc(
			sub,
			fn,
			proc,
		), nil
	}
	return proc
}
