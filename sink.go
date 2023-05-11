package sb

type Sink func(*Token) (Sink, error)

var _ SBUnmarshaler = Sink(nil)

func (s Sink) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if s == nil {
			return cont.Sink(token)
		}
		next, err := s(token)
		if err != nil { // NOCOVER
			return nil, err
		}
		return next.UnmarshalSB(ctx, cont), nil
	}
}

func (s Sink) Sink(token *Token) (Sink, error) {
	if s == nil {
		return nil, nil
	}
	return s(token)
}

func (s Sink) Marshal(o any) (Sink, error) {
	proc := Marshal(o)
	for {
		var token Token
		err := proc.Next(&token)
		if err != nil {
			return nil, err
		}
		if token.Invalid() {
			break
		}
		s, err = s.Sink(&token)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}
