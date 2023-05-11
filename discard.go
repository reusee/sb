package sb

func Discard(token *Token) (Sink, error) {
	if token.Invalid() {
		return nil, nil
	}
	return Discard, nil
}
