package sb

import (
	"io"
	"io/ioutil"
	"reflect"
)

type ChunkedReader struct {
	R io.Reader
	N int64
}

var _ SBMarshaler = new(ChunkedReader)

func (c ChunkedReader) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
				Kind: KindArray,
			}, c.marshal(
				ctx,
				func() (*Token, Proc, error) {
					return &Token{
						Kind: KindArrayEnd,
					}, cont, nil
				},
			), nil
	}
}

func (c ChunkedReader) marshal(ctx Ctx, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		bs, err := ioutil.ReadAll(
			&io.LimitedReader{
				R: c.R,
				N: c.N,
			},
		)
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if len(bs) > 0 {
			return nil, ctx.Marshal(
				ctx,
				reflect.ValueOf(bs),
				proc,
			), nil
		}
		return nil, cont, nil
	}
	return proc
}
