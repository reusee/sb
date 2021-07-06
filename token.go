package sb

type Token struct {
	Value any
	Kind  Kind
}

var _ SBMarshaler = Token{}

func (t Token) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &t, cont, nil
	}
}

var _ SBUnmarshaler = new(Token)

func (t *Token) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		*t = *token
		return cont, nil
	}
}

var Min = &Token{
	Kind: KindMin,
}

var Max = &Token{
	Kind: KindMax,
}

var NaN = &Token{
	Kind: KindNaN,
}

var Nil = &Token{
	Kind: KindNil,
}
