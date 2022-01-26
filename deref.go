package sb

func Deref(
	stream Stream,
	getStream func([]byte) (Stream, error),
) *Proc {
	proc := deref(
		stream,
		getStream,
		nil,
	)
	return &proc
}

func deref(
	stream Stream,
	getStream func([]byte) (Stream, error),
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
		if token.Kind == KindRef {
			subStream, err := getStream(token.Value.([]byte))
			if err != nil {
				return nil, nil, err
			}
			if subStream == nil {
				return nil, func() (*Token, Proc, error) {
					return token, proc, nil
				}, nil
			}
			return nil, IterStream(
				subStream,
				proc,
			), nil
		}
		return token, proc, nil
	}
	return proc
}
