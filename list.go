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

func (l *list) Next() (ret *Token) {
	if l.index >= len(l.tokens) {
		return nil
	}
	ret = &l.tokens[l.index]
	l.index++
	return
}
