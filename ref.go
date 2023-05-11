package sb

import (
	"io"
	"reflect"

	"github.com/reusee/e5"
)

type Ref []byte

var _ SBMarshaler = Ref{}

func (r Ref) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func(token *Token) (Proc, error) {
		token.Kind = KindRef
		token.Value = []byte(r)
		return cont, nil
	}
}

var _ SBUnmarshaler = new(Ref)

func (r *Ref) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(token *Token) (Sink, error) {
		if token.Invalid() {
			return nil, we.With(WithPath(ctx), e5.With(io.ErrUnexpectedEOF))(UnmarshalError)
		}
		if token.Kind != KindRef {
			return nil, we.With(WithPath(ctx), e5.With(TypeMismatch(token.Kind, reflect.Slice)))(UnmarshalError)
		}
		*r = Ref(token.Value.([]byte))
		return cont, nil
	}
}
