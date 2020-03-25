package sb

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
