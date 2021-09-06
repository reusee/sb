package sb

func ConcatSinks(sinks ...Sink) Sink {
	if len(sinks) == 0 {
		return nil
	}
	for sinks[0] == nil {
		sinks = sinks[1:]
		if len(sinks) == 0 {
			return nil
		}
	}
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		var err error
		sinks[0], err = sinks[0].Sink(token)
		if err != nil {
			return nil, err
		}
		for sinks[0] == nil {
			sinks = sinks[1:]
			if len(sinks) == 0 {
				return nil, nil
			}
		}
		return sink, nil
	}
	return sink
}

func ConcatProcs(procs ...Proc) Proc {
	if len(procs) == 0 {
		return nil
	}
	for procs[0] == nil {
		procs = procs[1:]
		if len(procs) == 0 {
			return nil
		}
	}
	var proc Proc
	proc = func() (*Token, Proc, error) {
		token, err := procs[0].Next()
		if err != nil {
			return nil, nil, err
		}
		if token == nil {
			procs = procs[1:]
			if len(procs) == 0 {
				return nil, nil, nil
			}
			return nil, proc, nil
		}
		return token, proc, nil
	}
	return proc
}
