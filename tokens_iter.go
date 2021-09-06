package sb

func (t Tokens) Iter() Proc {
	return IterTokens(t, 0, nil)
}

func IterTokens(
	tokens Tokens,
	index int,
	cont Proc,
) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if index >= len(tokens) {
			return nil, cont, nil
		}
		token := tokens[index]
		index++
		return &token, proc, nil
	}
	return proc
}
