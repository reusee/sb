package sb

func (t Tokens) Iter() *Proc {
	proc := IterTokens(t, 0, nil)
	return &proc
}

func IterTokens(
	tokens Tokens,
	index int,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		if index >= len(tokens) {
			return nil, cont, nil
		}
		return &tokens[index], IterTokens(tokens, index+1, cont), nil
	}
}
