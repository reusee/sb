package sb

import (
	"io"
	"reflect"

	"github.com/reusee/e4"
)

type Ref []byte

var _ SBMarshaler = Ref{}

func (r Ref) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return &Token{
			Kind:  KindRef,
			Value: []byte(r),
		}, cont, nil
	}
}

var _ SBUnmarshaler = new(Ref)

func (r *Ref) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, we.With(WithPath(ctx), e4.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}
		if token.Kind != KindRef {
			return nil, we.With(WithPath(ctx), e4.With(TypeMismatch(token.Kind, reflect.Slice)))(UnmarshalError)
		}
		*r = Ref(token.Value.([]byte))
		return cont, nil
	}
}
