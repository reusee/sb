package sb

type Proc func() (*Token, Proc, error)

var _ Stream = (*Proc)(nil)

func (p *Proc) Next() (*Token, error) {
	for {
		if *p == nil {
			return nil, nil
		}
		var ret *Token
		var err error
		ret, *p, err = (*p)()
		if ret != nil || err != nil {
			return ret, err
		}
	}
}

var _ SBMarshaler = Proc(nil)

func (p Proc) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		if p == nil {
			return nil, cont, nil
		}
		token, next, err := p()
		if err != nil {
			return nil, nil, err
		}
		return token, next.MarshalSB(ctx, cont), nil
	}
}
