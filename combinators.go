package sb

type Sink func(*Token) (Sink, error)

func AltSink(sinks ...Sink) Sink {
	return func(token *Token) (Sink, error) {
		next := make([]Sink, 0, len(sinks))
		var err error
		for _, sink := range sinks {
			sink, err = sink(token)
			if err != nil {
				continue
			}
			if sink == nil {
				return nil, nil
			}
			next = append(next, sink)
		}
		if len(next) == 0 {
			return nil, err
		}
		if len(next) == 1 {
			return next[0], nil
		}
		return AltSink(next...), nil
	}
}

func Copy(stream Stream, sink Sink) error {
	var token *Token
	var err error
	for {
		if stream != nil {
			token, err = stream.Next()
			if err != nil {
				return err
			}
			if token == nil {
				stream = nil
			}
		}
		if sink != nil {
			sink, err = sink(token)
			if err != nil {
				return err
			}
		}
		if sink == nil && stream == nil {
			break
		}
	}
	return nil
}

func ExpectKind(kind Kind, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil || token.Kind != kind {
			if err, ok := kindToExpectingErr[kind]; ok {
				return nil, UnmarshalError{err}
			} else {
				return nil, UnmarshalError{ExpectingValue}
			}
		}
		return cont, nil
	}
}

func Tee(stream Stream, sinks ...Sink) *Proc {
	proc := TeeProc(stream, sinks, nil)
	return &proc
}

func TeeProc(stream Stream, sinks []Sink, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		var token *Token
		var err error
		if stream != nil {
			token, err = stream.Next()
			if err != nil { // NOCOVER
				return nil, nil, err
			}
		}
		for i := 0; i < len(sinks); {
			sink, err := sinks[i](token)
			if err != nil { // NOCOVER
				return nil, nil, err
			}
			if sink == nil {
				sinks[i] = sinks[len(sinks)-1]
				sinks = sinks[:len(sinks)-1]
			} else {
				sinks[i] = sink
				i++
			}
		}
		if token == nil && len(sinks) == 0 {
			return nil, cont, nil
		}
		return token, TeeProc(stream, sinks, cont), nil
	}
}

func Discard(token *Token) (Sink, error) {
	if token == nil {
		return nil, nil
	}
	return Discard, nil
}

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

func FilterSink(sink Sink, fn func(*Token) bool) Sink {
	if sink == nil {
		return nil
	}
	return func(token *Token) (Sink, error) {
		var err error
		if token == nil || !fn(token) {
			sink, err = sink(token)
			if err != nil {
				return nil, err
			}
		}
		return FilterSink(sink, fn), nil
	}
}
