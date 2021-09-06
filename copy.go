package sb

func Copy(proc Proc, sinks ...Sink) error {
	var err error
	for {

		if len(sinks) == 0 {
			break
		}

		var token *Token
		for token == nil {
			if proc != nil {
				token, err = proc.Next()
				if err != nil {
					return err
				}
				if token == nil {
					proc = nil
				}
			} else {
				break
			}
		}

		if len(sinks) > 0 {
			for i := 0; i < len(sinks); {
				sink := sinks[i]
				if sink == nil {
					sinks[i] = sinks[len(sinks)-1]
					sinks = sinks[:len(sinks)-1]
					continue
				}
				sink, err = sink(token)
				if err != nil {
					return err
				}
				if sink == nil {
					sinks[i] = sinks[len(sinks)-1]
					sinks = sinks[:len(sinks)-1]
					continue
				}
				sinks[i] = sink
				i++
			}
		}

		if len(sinks) == 0 && proc == nil {
			break
		}
	}
	return nil
}
