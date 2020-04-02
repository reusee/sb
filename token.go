package sb

type Token struct {
	Kind  Kind
	Value any
}

type min struct{}

var _ SBMarshaler = min{}

func (_ min) MarshalSB(_ Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindMin,
		}, cont, nil
	}
}

var Min = min{}

type max struct{}

var _ SBMarshaler = max{}

func (_ max) MarshalSB(_ Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindMax,
		}, cont, nil
	}
}

var Max = max{}
