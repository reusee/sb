package sb

func FilterProc(
	proc Proc,
	predict func(*Token) bool,
) Proc {
	return filterProc(proc, predict, nil)
}

func filterProc(
	src Proc,
	predict func(*Token) bool,
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
		if !predict(token) {
			token = nil
		}
		return token, proc, nil
	}
	return proc
}

func FilterSink(sink Sink, fn func(*Token) bool) Sink {
	var s Sink
	s = func(token *Token) (Sink, error) {
		var err error
		if token == nil || fn(token) {
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
