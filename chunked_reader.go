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
	return MarshalReaderChunked(ctx, c.R, c.N, cont)
}

func MarshalReaderChunked(ctx Ctx, r io.Reader, n int64, cont Proc) Proc {
	var marshal Proc
	marshal = func() (*Token, Proc, error) {
		bs, err := ioutil.ReadAll(
			&io.LimitedReader{
				R: r,
				N: n,
			},
		)
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if len(bs) > 0 {
			return nil, ctx.Marshal(
				ctx,
				reflect.ValueOf(bs),
				marshal,
			), nil
		}
		return &Token{
			Kind: KindArrayEnd,
		}, cont, nil
	}
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindArray,
		}, marshal, nil
	}
}
