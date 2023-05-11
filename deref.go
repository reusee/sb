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
	proc = func(token *Token) (Proc, error) {
		err := stream.Next(token)
		if err != nil { // NOCOVER
			return nil, err
		}
		if token.Invalid() {
			return cont, nil
		}
		if token.Kind == KindRef {
			subStream, err := getStream(token.Value.([]byte))
			if err != nil {
				return nil, err
			}
			if subStream == nil {
				return func(t *Token) (Proc, error) {
					*t = *token
					return proc, nil
				}, nil
			}
			token.Reset() // do not provide KindRef token
			return IterStream(
				subStream,
				proc,
			), nil
		}
		return proc, nil
	}
	return proc
}
