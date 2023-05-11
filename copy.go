package sb

func Copy(stream Stream, sinks ...Sink) error {
	var err error
	for len(sinks) > 0 {

		var token Token
		for !token.Valid() {
			if stream != nil {
				err = stream.Next(&token)
				if err != nil {
					return err
				}
				if !token.Valid() {
					stream = nil
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
				sink, err = sink(&token)
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

		if len(sinks) == 0 && stream == nil {
			break
		}
	}
	return nil
}
