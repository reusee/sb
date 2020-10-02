package sb

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

func Hash(
	newState func() hash.Hash,
	sum *[]byte,
	cont Sink,
) Sink {
	return HashFunc(
		newState, sum, nil, cont,
	)
}

func HashFunc(
	newState func() hash.Hash,
	sum *[]byte,
	fn func([]byte, *Token) error,
	cont Sink,
) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, ExpectingValue
		}

		state := newState()
		if err := binary.Write(
			state,
			binary.LittleEndian,
			token.Kind,
		); err != nil { // NOCOVER
			return nil, err
		}

		if fn != nil {
			if err := fn(nil, token); err != nil {
				return nil, err
			}
		}

		switch token.Kind {

		case KindInvalid,
			KindMin,
			KindNil,
			KindNaN,
			KindMax,
			KindArrayEnd,
			KindObjectEnd,
			KindMapEnd,
			KindTupleEnd:
			*sum = state.Sum(nil)
			if fn != nil {
				if err := fn(*sum, token); err != nil {
					return nil, err
				}
			}
			return cont, nil

		case KindBool,
			KindString,
			KindBytes,
			KindInt,
			KindInt8,
			KindInt16,
			KindInt32,
			KindInt64,
			KindUint,
			KindUint8,
			KindUint16,
			KindUint32,
			KindUint64,
			KindFloat32,
			KindFloat64:
			if token.Kind == KindInt {
				if err := binary.Write(state, binary.LittleEndian, int64(token.Value.(int))); err != nil { // NOCOVER
					return nil, err
				}
			} else if token.Kind == KindUint {
				if err := binary.Write(state, binary.LittleEndian, uint64(token.Value.(uint))); err != nil { // NOCOVER
					return nil, err
				}
			} else if token.Kind == KindString {
				if _, err := io.WriteString(state, token.Value.(string)); err != nil { // NOCOVER
					return nil, err
				}
			} else {
				if err := binary.Write(state, binary.LittleEndian, token.Value); err != nil { // NOCOVER
					return nil, err
				}
			}
			*sum = state.Sum(nil)
			if fn != nil {
				if err := fn(*sum, token); err != nil {
					return nil, err
				}
			}
			return cont, nil

		case KindArray, KindObject, KindMap, KindTuple:
			t := token
			return HashCompound(
				newState,
				state,
				fn,
				func(token *Token) (Sink, error) {
					*sum = state.Sum(nil)
					if fn != nil {
						if err := fn(*sum, t); err != nil {
							return nil, err
						}
					}
					return cont.Sink(token)
				},
			), nil

		default: // NOCOVER
			panic(fmt.Errorf("unexpected token: %+v", token))

		}

	}
}

func HashCompound(
	newState func() hash.Hash,
	state hash.Hash,
	fn func([]byte, *Token) error,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return nil, ExpectingValue
		}

		var subHash []byte
		var next Sink

		if token.Kind == KindArrayEnd ||
			token.Kind == KindObjectEnd ||
			token.Kind == KindMapEnd ||
			token.Kind == KindTupleEnd {
			next = cont
		} else {
			next = sink
		}

		return HashFunc(
			newState,
			&subHash,
			fn,
			func(token *Token) (Sink, error) {
				if _, err := state.Write(subHash); err != nil { // NOCOVER
					return nil, err
				}
				if next != nil {
					return next(token)
				}
				return nil, nil // NOCOVER
			},
		)(token)

	}

	return sink
}
