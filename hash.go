package sb

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math"
)

func Hash(
	newState func() hash.Hash,
	target *[]byte,
	cont Sink,
) Sink {
	return HashFunc(
		newState, target, nil, cont,
	)
}

func HashFunc(
	newState func() hash.Hash,
	target *[]byte,
	fn func([]byte, *Token) error,
	cont Sink,
) Sink {
	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, ExpectingValue
		}

		state := newState()
		if _, err := state.Write([]byte{
			byte(token.Kind),
		}); err != nil {
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
			sum := state.Sum(nil)
			if target != nil {
				*target = sum
			}
			if fn != nil {
				if err := fn(sum, token); err != nil {
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

			buf, put := getEightBytes()
			defer put()

			switch token.Kind {
			case KindBool:
				if token.Value.(bool) {
					if _, err := state.Write([]byte{1}); err != nil {
						return nil, err
					}
				} else {
					if _, err := state.Write([]byte{0}); err != nil {
						return nil, err
					}
				}
			case KindString:
				if _, err := io.WriteString(state, token.Value.(string)); err != nil { // NOCOVER
					return nil, err
				}
			case KindBytes:
				if _, err := state.Write(token.Value.([]byte)); err != nil {
					return nil, err
				}
			case KindInt:
				binary.LittleEndian.PutUint64(buf, uint64(token.Value.(int)))
				if _, err := state.Write(buf); err != nil {
					return nil, err
				}
			case KindInt8:
				if _, err := state.Write([]byte{uint8(token.Value.(int8))}); err != nil {
					return nil, err
				}
			case KindInt16:
				binary.LittleEndian.PutUint16(buf, uint16(token.Value.(int16)))
				if _, err := state.Write((buf)[:2]); err != nil {
					return nil, err
				}
			case KindInt32:
				binary.LittleEndian.PutUint32(buf, uint32(token.Value.(int32)))
				if _, err := state.Write((buf)[:4]); err != nil {
					return nil, err
				}
			case KindInt64:
				binary.LittleEndian.PutUint64(buf, uint64(token.Value.(int64)))
				if _, err := state.Write((buf)); err != nil {
					return nil, err
				}
			case KindUint:
				binary.LittleEndian.PutUint64(buf, uint64(token.Value.(uint)))
				if _, err := state.Write(buf); err != nil {
					return nil, err
				}
			case KindUint8:
				if _, err := state.Write([]byte{token.Value.(uint8)}); err != nil {
					return nil, err
				}
			case KindUint16:
				binary.LittleEndian.PutUint16(buf, token.Value.(uint16))
				if _, err := state.Write((buf)[:2]); err != nil {
					return nil, err
				}
			case KindUint32:
				binary.LittleEndian.PutUint32(buf, token.Value.(uint32))
				if _, err := state.Write((buf)[:4]); err != nil {
					return nil, err
				}
			case KindUint64:
				binary.LittleEndian.PutUint64(buf, token.Value.(uint64))
				if _, err := state.Write((buf)); err != nil {
					return nil, err
				}
			case KindFloat32:
				binary.LittleEndian.PutUint32(buf, math.Float32bits(token.Value.(float32)))
				if _, err := state.Write((buf)[:4]); err != nil {
					return nil, err
				}
			case KindFloat64:
				binary.LittleEndian.PutUint64(buf, math.Float64bits(token.Value.(float64)))
				if _, err := state.Write((buf)); err != nil {
					return nil, err
				}
			default:
				panic("impossible")
			}
			sum := state.Sum(nil)
			if target != nil {
				*target = sum
			}
			if fn != nil {
				if err := fn(sum, token); err != nil {
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
					sum := state.Sum(nil)
					if target != nil {
						*target = sum
					}
					if fn != nil {
						if err := fn(sum, t); err != nil {
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
