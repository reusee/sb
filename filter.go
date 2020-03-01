package sb

func FilterProc(
	stream Stream,
	predict func(*Token) bool,
) *Proc {
	proc := filterProc(stream, predict, nil)
	return &proc
}

func filterProc(
	stream Stream,
	predict func(*Token) bool,
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
		if predict(token) {
			token = nil
		}
		return token, filterProc(stream, predict, cont), nil
	}
}
