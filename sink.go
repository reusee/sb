package sb

type Sink func(*Token) (Sink, error)

var _ SBUnmarshaler = Sink(nil)

func (s Sink) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if s == nil {
			return cont, nil
		}
		next, err := s(token)
		if err != nil { // NOCOVER
			return nil, err
		}
		return next.UnmarshalSB(ctx, cont), nil
	}
}
