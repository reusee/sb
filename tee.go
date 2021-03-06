package sb

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
