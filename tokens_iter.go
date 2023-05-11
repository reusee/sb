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
	proc = func(token *Token) (Proc, error) {
		if index >= len(tokens) {
			return cont, nil
		}
		*token = tokens[index]
		index++
		return proc, nil
	}
	return proc
}
