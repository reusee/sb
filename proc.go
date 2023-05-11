package sb

type Proc func(token *Token) (Proc, error)

type Stream = *Proc

func (p *Proc) Next(token *Token) (err error) {
	for !token.Valid() {
		if p != nil && *p != nil {
			*p, err = (*p)(token)
			if err != nil {
				return
			}
		} else {
			break
		}
	}
	return
}

var _ SBMarshaler = Proc(nil)

func (p Proc) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func(token *Token) (Proc, error) {
		if p == nil {
			return cont, nil
		}
		next, err := p(token)
		if err != nil { // NOCOVER
			return nil, err
		}
		return next.MarshalSB(ctx, cont), nil
	}
}
