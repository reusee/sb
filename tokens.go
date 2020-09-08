package sb

type Tokens []Token

func TokensFromStream(stream Stream) (tokens Tokens, err error) {
	for {
		p, err := stream.Next()
		if err != nil {
			return nil, err
		}
		if p == nil {
			break
		}
		tokens = append(tokens, *p)
	}
	return
}

func MustTokensFromStream(stream Stream) Tokens {
	tokens, err := TokensFromStream(stream)
	if err != nil {
		panic(err)
	}
	return tokens
}

func CollectTokens(tokens *Tokens) Sink {
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return nil, nil
		}
		*tokens = append(*tokens, *token)
		return sink, nil
	}
	return sink
}

func CollectValueTokens(tokens *Tokens) Sink {
	var sink Sink
	var stack []Kind
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			if len(stack) > 0 {
				return nil, ExpectingValue
			}
			return nil, nil // NOCOVER
		}
		*tokens = append(*tokens, *token)
		switch token.Kind {
		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			if len(stack) == 0 {
				return nil, UnexpectedEndToken
			}
			if token.Kind != stack[len(stack)-1] {
				return nil, kindToExpectingErr[stack[len(stack)-1]]
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return nil, nil
			}
			return sink, nil
		case KindArray:
			stack = append(stack, KindArrayEnd)
			return sink, nil
		case KindObject:
			stack = append(stack, KindObjectEnd)
			return sink, nil
		case KindMap:
			stack = append(stack, KindMapEnd)
			return sink, nil
		case KindTuple:
			stack = append(stack, KindTupleEnd)
			return sink, nil
		}
		if len(stack) > 0 {
			return sink, nil
		}
		return nil, nil
	}
	return sink
}
