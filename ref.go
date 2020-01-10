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
	return func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil {
			return nil, nil, err
		}
		if token == nil {
			return nil, cont, nil
		}
		if token.Kind == KindRef {
			hash := token.Value.([]byte)
			subStream, err := getStream(hash)
			if err != nil {
				return nil, nil, err
			}
			return nil, IterStream(
				subStream,
				deref(stream, getStream, cont),
			), nil
		}
		return token, deref(stream, getStream, cont), nil
	}
}
