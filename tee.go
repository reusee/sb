package sb

func Tee(stream Stream, sinks ...Sink) *Proc {
	proc := TeeProc(stream, sinks, nil)
	return &proc
}

func TeeProc(stream Stream, sinks []Sink, cont Proc) Proc {
	var proc Proc
	proc = func(token *Token) (Proc, error) {
		var err error
		if stream != nil {
			err = stream.Next(token)
			if err != nil { // NOCOVER
				return nil, err
			}
		}
		for i := 0; i < len(sinks); {
			sink, err := sinks[i](token)
			if err != nil { // NOCOVER
				return nil, err
			}
			if sink == nil {
				sinks[i] = sinks[len(sinks)-1]
				sinks = sinks[:len(sinks)-1]
			} else {
				sinks[i] = sink
				i++
			}
		}
		if token.Invalid() && len(sinks) == 0 {
			return cont, nil
		}
		return proc, nil
	}
	return proc
}
