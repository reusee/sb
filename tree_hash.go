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

	if _, ok := t.Tags.Get("hash"); ok {
		return
	}

	if t.Token == nil {
		panic("empty tree")
	}
	token := t.Token

	if token.Kind == KindPostTag {
		panic("unexpected KindPostTag token")
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
		t.Tags.Set("hash", state.Sum(nil))

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
		t.Tags.Set("hash", state.Sum(nil))

	case KindArray, KindObject, KindMap, KindTuple:
		// write sub hashes
		for _, sub := range t.Subs {
			if err = sub.FillHash(newState); err != nil { // NOCOVER
				return
			}
			subHash, ok := sub.Tags.Get("hash")
			if !ok {
				panic("impossible")
			}
			if _, err = state.Write(subHash); err != nil { // NOCOVER
				return
			}
		}
		t.Tags.Set("hash", state.Sum(nil))

	default:
		panic(fmt.Errorf("unexpected token: %+v", token))

	}

	return
}
