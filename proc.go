package sb

type Proc func() (*Token, Proc, error)

func (p *Proc) Next() (ret *Token, err error) {
	for ret == nil {
		if p != nil && *p != nil {
			ret, *p, err = (*p)()
		} else {
			break
		}
	}
	return
}

var _ SBMarshaler = Proc(nil)

func (p Proc) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		if p == nil {
			return nil, cont, nil
		}
		token, next, err := p()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		return token, next.MarshalSB(ctx, cont), nil
	}
}
