package sb

func IterStream(
	stream Stream,
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
		return proc, nil
	}
	return proc
}
