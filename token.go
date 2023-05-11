package sb

type Token struct {
	Value any
	Kind  Kind
}

func (t *Token) Valid() bool {
	return t.Kind != KindInvalid
}

func (t *Token) Invalid() bool {
	return t.Kind == KindInvalid
}

func (t *Token) Reset() {
	*t = Token{}
}

var _ SBMarshaler = Token{}

func (t Token) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func(token *Token) (Proc, error) {
		*token = t
		return cont, nil
	}
}

var _ SBUnmarshaler = new(Token)

func (t *Token) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		*t = *token
		return cont, nil
	}
}

var Min = Token{
	Kind: KindMin,
}

var Max = Token{
	Kind: KindMax,
}

var NaN = Token{
	Kind: KindNaN,
}

var Nil = Token{
	Kind: KindNil,
}
