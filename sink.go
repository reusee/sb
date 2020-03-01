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

func Pipe(stream Stream, sinks ...Sink) error {
	var err error
	for {
		token, err := stream.Next()
		if err != nil {
			return err
		}
		if token == nil {
			break
		}
		if len(sinks) == 0 {
			return UnmarshalError{EmptySink}
		}
		for i := 0; i < len(sinks); {
			if sinks[i] == nil {
				sinks[i] = sinks[len(sinks)-1]
				sinks = sinks[:len(sinks)-1]
				continue
			}
			sinks[i], err = sinks[i](token)
			if err != nil {
				return err
			}
			i++
		}
	}

	for len(sinks) > 0 {
		for i := 0; i < len(sinks); {
			if sinks[i] == nil {
				sinks[i] = sinks[len(sinks)-1]
				sinks = sinks[:len(sinks)-1]
				continue
			}
			sinks[i], err = sinks[i](nil)
			if err != nil {
				return err
			}
			i++
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
