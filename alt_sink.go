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
