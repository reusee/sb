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
