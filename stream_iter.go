package sb

func IterStream(
	stream Stream,
	cont Proc,
) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token == nil {
			return nil, cont, nil
		}
		return token, proc, nil
	}
	return proc
}
