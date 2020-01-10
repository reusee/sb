package sb

func IterStream(
	stream Stream,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil {
			return nil, nil, err
		}
		if token == nil {
			return nil, cont, nil
		}
		return token, IterStream(stream, cont), nil
	}
}
