package sb

func Filter(
	stream Stream,
	predict func(*Token) bool,
) *Proc {
	proc := filter(stream, predict, nil)
	return &proc
}

func filter(
	stream Stream,
	predict func(*Token) bool,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil {
			return nil, nil, err
		}
		if predict(token) {
			return nil, filter(stream, predict, cont), nil
		}
		return token, filter(stream, predict, cont), nil
	}
}
