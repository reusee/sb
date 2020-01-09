package sb

import (
	"encoding/binary"
	"hash"
	"io"
)

type HashProc func() (*Token, HashProc, error)

func NewHasher(
	stream Stream,
	newState func() hash.Hash,
) *HashProc {
	states := []hash.Hash{
		newState(),
	}
	hasher := HashStream(
		stream,
		newState,
		&states,
		func() (*Token, HashProc, error) {
			if len(states) != 1 { // NOCOVER
				panic("bad state")
			}
			return &Token{
				Kind:  KindHash,
				Value: states[0].Sum(nil),
			}, nil, nil
		},
	)
	return &hasher
}

func HashStream(
	stream Stream,
	newState func() hash.Hash,
	states *[]hash.Hash,
	cont HashProc,
) HashProc {
	return func() (*Token, HashProc, error) {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			return nil, nil, err
		}
		if token == nil {
			// stop
			return nil, cont, nil
		}
		if token.Kind == KindHash {
			// rip hash tokens
			return nil, HashStream(stream, newState, states, cont), nil
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
				HashStream(stream, newState, states, cont),
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
				HashStream(stream, newState, states, cont),
			), nil

		case KindArray, KindObject, KindMap, KindTuple:
			// push state
			*states = append(*states, state)
			return token, HashStream(stream, newState, states, cont), nil

		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			// pop state
			sum := (*states)[len(*states)-1].Sum(nil)
			*states = (*states)[:len(*states)-1]
			return token, emitHash(
				sum, states,
				HashStream(stream, newState, states, cont),
			), nil

		}

		panic("impossible") // NOCOVER
	}
}

func emitHash(sum []byte, states *[]hash.Hash, cont HashProc) HashProc {
	return func() (*Token, HashProc, error) {
		// write to stack
		if _, err := (*states)[len(*states)-1].Write(sum); err != nil { // NOCOVER
			return nil, nil, err
		}
		return &Token{
			Kind:  KindHash,
			Value: sum,
		}, cont, nil
	}
}

var _ Stream = (*HashProc)(nil)

func (p *HashProc) Next() (*Token, error) {
	for {
		if p == nil || *p == nil {
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

func HashSum(
	stream Stream,
	newState func() hash.Hash,
) (
	sum []byte,
	err error,
) {
	hasher := NewHasher(stream, newState)
	var token, last *Token
	for {
		token, err = hasher.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			if last.Kind != KindHash {
				panic("bad hasher")
			}
			sum = last.Value.([]byte)
			break
		}
		last = token
	}
	return
}
