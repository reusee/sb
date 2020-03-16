package sb

func AltSink(sinks ...Sink) Sink {
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		var err error
		for i := 0; i < len(sinks); {
			sink := sinks[i]
			sink, err = sink(token)
			if err != nil {
				sinks[i] = sinks[len(sinks)-1]
				sinks = sinks[:len(sinks)-1]
				continue
			}
			if sink == nil {
				return nil, nil
			}
			sinks[i] = sink
			i++
		}
		if len(sinks) == 0 {
			return nil, err
		}
		if len(sinks) == 1 {
			return sinks[0], nil
		}
		return sink, nil
	}
	return sink
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
	var proc Proc
	proc = func() (*Token, Proc, error) {
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
		return token, proc, nil
	}
	return proc
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
	var proc Proc
	proc = func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token == nil {
			return nil, cont, nil
		}
		if predict(token) {
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
		if token == nil || !fn(token) {
			if sink == nil {
				return nil, nil
			}
			sink, err = sink(token)
			if err != nil {
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
