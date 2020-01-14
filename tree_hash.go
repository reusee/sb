package sb

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

func TreeHashSum(
	tree *Tree,
	newState func() hash.Hash,
) (
	sum []byte,
	err error,
) {

	if len(tree.Hash) > 0 {
		sum = tree.Hash
		return
	}

	if tree.Token == nil {
		panic("empty tree")
	}
	token := tree.Token

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
		sum = state.Sum(nil)
		return

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
		sum = state.Sum(nil)
		return

	case KindArray, KindObject, KindMap, KindTuple:
		// write sub hashes
		for _, sub := range tree.Subs {
			var subHash []byte
			subHash, err = TreeHashSum(sub, newState)
			if err != nil {
				return
			}
			if _, err = state.Write(subHash); err != nil {
				return
			}
		}
		sum = state.Sum(nil)
		return

	default:
		panic(fmt.Errorf("unexpected token: %+v", token))

	}

	return
}
