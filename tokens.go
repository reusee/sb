package sb

type Tokens []Token

func TokensFromProc(proc Proc) (tokens Tokens, err error) {
	for {
		p, err := proc.Next()
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

func MustTokensFromProc(proc Proc) Tokens {
	tokens, err := TokensFromProc(proc)
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
	type Frame struct {
		Kind Kind
		N    int
	}
	var stack []*Frame
	sink = func(token *Token) (Sink, error) {
		if len(stack) > 0 {
			parent := stack[len(stack)-1]
			if parent.Kind == KindTypeName && parent.N > 0 {
				stack = stack[:len(stack)-1]
				if len(stack) > 0 {
					parent = stack[len(stack)-1]
				}
			}
			parent.N++
		}
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
			if token.Kind != stack[len(stack)-1].Kind {
				return nil, kindToExpectingErr[stack[len(stack)-1].Kind]
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return nil, nil
			}
			return sink, nil
		case KindArray:
			stack = append(stack, &Frame{
				Kind: KindArrayEnd,
			})
			return sink, nil
		case KindObject:
			stack = append(stack, &Frame{
				Kind: KindObjectEnd,
			})
			return sink, nil
		case KindMap:
			stack = append(stack, &Frame{
				Kind: KindMapEnd,
			})
			return sink, nil
		case KindTuple:
			stack = append(stack, &Frame{
				Kind: KindTupleEnd,
			})
			return sink, nil
		case KindTypeName:
			stack = append(stack, &Frame{
				Kind: KindTypeName,
			})
			return sink, nil
		}
		if len(stack) > 0 {
			return sink, nil
		}
		return nil, nil
	}
	return sink
}
