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
	var proc Proc
	proc = func(token *Token) (Proc, error) {
		err := stream.Next(token)
		if err != nil { // NOCOVER
			return nil, err
		}
		if token.Invalid() {
			return cont, nil
		}
		if !predict(token) {
			token.Reset()
		}
		return proc, nil
	}
	return proc
}

func FilterSink(sink Sink, fn func(*Token) bool) Sink {
	var s Sink
	s = func(token *Token) (Sink, error) {
		var err error
		if token.Invalid() || fn(token) {
			if sink == nil {
				return nil, nil
			}
			sink, err = sink(token)
			if err != nil { // NOCOVER
				return nil, err
			}
		}
		if sink == nil {
			return nil, nil
		}
		return s, nil
	}
	return s
}
