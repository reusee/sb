package sb

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

func (t *Tree) FillHash(
	newState func() hash.Hash,
) (
	err error,
) {

	if len(t.Hash) > 0 {
		return
	}

	if t.Token == nil {
		panic("empty tree")
	}
	token := t.Token

	if token.Kind == KindPostHash {
		panic("unexpected KindPostHash token")
	}

	state := newState()
	if err = binary.Write(state, binary.LittleEndian, token.Kind); err != nil { // NOCOVER
		return
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
		t.Hash = state.Sum(nil)

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
			if err = binary.Write(state, binary.LittleEndian, int64(token.Value.(int))); err != nil { // NOCOVER
				return
			}
		} else if token.Kind == KindUint {
			if err = binary.Write(state, binary.LittleEndian, uint64(token.Value.(uint))); err != nil { // NOCOVER
				return
			}
		} else if token.Kind == KindString {
			if _, err = io.WriteString(state, token.Value.(string)); err != nil { // NOCOVER
				return
			}
		} else {
			if err = binary.Write(state, binary.LittleEndian, token.Value); err != nil { // NOCOVER
				return
			}
		}
		t.Hash = state.Sum(nil)

	case KindArray, KindObject, KindMap, KindTuple:
		// write sub hashes
		for _, sub := range t.Subs {
			if err = sub.FillHash(newState); err != nil {
				return
			}
			if _, err = state.Write(sub.Hash); err != nil {
				return
			}
		}
		t.Hash = state.Sum(nil)

	default:
		panic(fmt.Errorf("unexpected token: %+v", token))

	}

	return
}
