package sb

type list struct {
	tokens []Token
	index  int
}

func List(tokens []Token) *list {
	return &list{
		tokens: tokens,
		index:  0,
	}
}

var _ Stream = new(list)

func (l *list) Next() (ret *Token, err error) {
	if l.index >= len(l.tokens) {
		return nil, nil
	}
	ret = &l.tokens[l.index]
	l.index++
	return
}

func (l *list) Peek() (ret *Token, err error) {
	if l.index >= len(l.tokens) {
		return nil, nil
	}
	ret = &l.tokens[l.index]
	return
}
