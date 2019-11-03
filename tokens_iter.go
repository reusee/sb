package sb

type TokensIter struct {
	tokens Tokens
	index  int
}

func (t Tokens) Iter() *TokensIter {
	return &TokensIter{
		tokens: t,
		index:  0,
	}
}

var _ Stream = new(TokensIter)

func (t *TokensIter) Next() (ret *Token, err error) {
	if t.index >= len(t.tokens) {
		return nil, nil
	}
	ret = &t.tokens[t.index]
	t.index++
	return
}
