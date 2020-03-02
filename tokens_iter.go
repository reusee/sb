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
