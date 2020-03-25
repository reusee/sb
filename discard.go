package sb

func Discard(token *Token) (Sink, error) {
	if token == nil {
		return nil, nil
	}
	return Discard, nil
}
