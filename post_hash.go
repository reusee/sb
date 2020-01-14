package sb

import (
	"encoding/binary"
	"hash"
	"io"
)

func NewPostHasher(
	stream Stream,
	newState func() hash.Hash,
) *Proc {
	hasher := PostHashStream(
		stream,
		newState,
		&[]hash.Hash{},
		nil,
	)
	return &hasher
}

func PostHashStream(
	stream Stream,
	newState func() hash.Hash,
	states *[]hash.Hash,
	cont Proc,
) Proc {
	return func() (*Token, Proc, error) {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token == nil {
			// stop
			return nil, cont, nil
		}
		if token.Kind == KindPostHash {
			// rip hash tokens
			return nil, PostHashStream(stream, newState, states, cont), nil
		}

		state := newState()
		if err := binary.Write(state, binary.LittleEndian, token.Kind); err != nil { // NOCOVER
			return nil, nil, err
		}

		switch token.Kind {

		case KindInvalid,
			KindMin,
			KindNil,
			KindNaN,
			KindMax:
			sum := state.Sum(nil)
			return token, emitHash(
				sum, states,
				PostHashStream(stream, newState, states, cont),
			), nil

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
					return nil, nil, err
				}
			} else if token.Kind == KindUint {
				if err := binary.Write(state, binary.LittleEndian, uint64(token.Value.(uint))); err != nil { // NOCOVER
					return nil, nil, err
				}
			} else if token.Kind == KindString {
				if _, err := io.WriteString(state, token.Value.(string)); err != nil { // NOCOVER
					return nil, nil, err
				}
			} else {
				if err := binary.Write(state, binary.LittleEndian, token.Value); err != nil { // NOCOVER
					return nil, nil, err
				}
			}
			sum := state.Sum(nil)
			return token, emitHash(
				sum, states,
				PostHashStream(stream, newState, states, cont),
			), nil

		case KindArray, KindObject, KindMap, KindTuple:
			// push state
			*states = append(*states, state)
			return token, PostHashStream(stream, newState, states, cont), nil

		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			// pop state
			sum := (*states)[len(*states)-1].Sum(nil)
			*states = (*states)[:len(*states)-1]
			return token, emitHash(
				sum, states,
				PostHashStream(stream, newState, states, cont),
			), nil

		}

		panic("impossible") // NOCOVER
	}
}

func emitHash(sum []byte, states *[]hash.Hash, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		// write to stack
		if len(*states) > 0 {
			if _, err := (*states)[len(*states)-1].Write(sum); err != nil { // NOCOVER
				return nil, nil, err
			}
		}
		return &Token{
			Kind:  KindPostHash,
			Value: sum,
		}, cont, nil
	}
}
