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
