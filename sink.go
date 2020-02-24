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
