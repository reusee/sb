package sb

func Deref(
	proc Proc,
	get func([]byte) (Proc, error),
) Proc {
	return deref(
		proc,
		get,
		nil,
	)
}

func deref(
	src Proc,
	get func([]byte) (Proc, error),
	cont Proc,
) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		token, err := src.Next()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token == nil {
			return nil, cont, nil
		}
		if token.Kind == KindRef {
			sub, err := get(token.Value.([]byte))
			if err != nil {
				return nil, nil, err
			}
			if sub == nil {
				return nil, func() (*Token, Proc, error) {
					return token, proc, nil
				}, nil
			}
			return nil, Iter(
				sub,
				proc,
			), nil
		}
		return token, proc, nil
	}
	return proc
}
