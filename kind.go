package sb

func ExpectKind(kind Kind, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil || token.Kind != kind {
			if err, ok := kindToExpectingErr[kind]; ok {
				return nil, UnmarshalError{err}
			} else {
				return nil, UnmarshalError{ExpectingValue}
			}
		}
		return cont, nil
	}
}
