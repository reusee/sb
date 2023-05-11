package sb

func (t *Tree) Iter() *Proc {
	proc := IterTree(t, nil)
	return &proc
}

func IterTree(
	tree *Tree,
	cont Proc,
) Proc {
	return func(token *Token) (Proc, error) {
		*token = *tree.Token
		return IterSubTrees(
			tree.Subs, 0,
			func(token *Token) (Proc, error) {
				return cont, nil
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
	proc = func(_ *Token) (Proc, error) {
		if len(subs) == 0 {
			return cont, nil
		}
		if index >= len(subs) {
			return cont, nil
		}
		sub := subs[index]
		index++
		return IterTree(
			sub,
			proc,
		), nil
	}
	return proc
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
	return func(token *Token) (Proc, error) {
		t, err := fn(tree)
		if err != nil { // NOCOVER
			return nil, err
		}
		if t != nil {
			*token = *t
			return cont, nil
		}
		*token = *tree.Token
		return IterSubTreesFunc(
			tree.Subs, 0, fn,
			func(token *Token) (Proc, error) {
				return cont, nil
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
	proc = func(_ *Token) (Proc, error) {
		if len(subs) == 0 {
			return cont, nil
		}
		if index >= len(subs) {
			return cont, nil
		}
		sub := subs[index]
		index++
		return IterTreeFunc(
			sub,
			fn,
			proc,
		), nil
	}
	return proc
}
