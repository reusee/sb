package sb

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math"
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

	state := newState()
	if _, err = state.Write([]byte{byte(token.Kind)}); err != nil {
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

		buf := eightBytesPool.Get().(*[]byte)
		defer eightBytesPool.Put(buf)

		switch token.Kind {
		case KindBool:
			if token.Value.(bool) {
				if _, err := state.Write([]byte{1}); err != nil {
					return err
				}
			} else {
				if _, err := state.Write([]byte{0}); err != nil {
					return err
				}
			}
		case KindString:
			if _, err := io.WriteString(state, token.Value.(string)); err != nil { // NOCOVER
				return err
			}
		case KindBytes:
			if _, err := state.Write(token.Value.([]byte)); err != nil {
				return err
			}
		case KindInt:
			binary.LittleEndian.PutUint64(*buf, uint64(token.Value.(int)))
			if _, err := state.Write(*buf); err != nil {
				return err
			}
		case KindInt8:
			if _, err := state.Write([]byte{uint8(token.Value.(int8))}); err != nil {
				return err
			}
		case KindInt16:
			binary.LittleEndian.PutUint16(*buf, uint16(token.Value.(int16)))
			if _, err := state.Write((*buf)[:2]); err != nil {
				return err
			}
		case KindInt32:
			binary.LittleEndian.PutUint32(*buf, uint32(token.Value.(int32)))
			if _, err := state.Write((*buf)[:4]); err != nil {
				return err
			}
		case KindInt64:
			binary.LittleEndian.PutUint64(*buf, uint64(token.Value.(int64)))
			if _, err := state.Write((*buf)); err != nil {
				return err
			}
		case KindUint:
			binary.LittleEndian.PutUint64(*buf, uint64(token.Value.(uint)))
			if _, err := state.Write(*buf); err != nil {
				return err
			}
		case KindUint8:
			if _, err := state.Write([]byte{token.Value.(uint8)}); err != nil {
				return err
			}
		case KindUint16:
			binary.LittleEndian.PutUint16(*buf, token.Value.(uint16))
			if _, err := state.Write((*buf)[:2]); err != nil {
				return err
			}
		case KindUint32:
			binary.LittleEndian.PutUint32(*buf, token.Value.(uint32))
			if _, err := state.Write((*buf)[:4]); err != nil {
				return err
			}
		case KindUint64:
			binary.LittleEndian.PutUint64(*buf, token.Value.(uint64))
			if _, err := state.Write((*buf)); err != nil {
				return err
			}
		case KindFloat32:
			binary.LittleEndian.PutUint32(*buf, math.Float32bits(token.Value.(float32)))
			if _, err := state.Write((*buf)[:4]); err != nil {
				return err
			}
		case KindFloat64:
			binary.LittleEndian.PutUint64(*buf, math.Float64bits(token.Value.(float64)))
			if _, err := state.Write((*buf)); err != nil {
				return err
			}
		default:
			panic("impossible")
		}
		t.Hash = state.Sum(nil)

	case KindArray, KindObject, KindMap, KindTuple:
		// write sub hashes
		for _, sub := range t.Subs {
			if err = sub.FillHash(newState); err != nil { // NOCOVER
				return
			}
			if _, err = state.Write(sub.Hash); err != nil { // NOCOVER
				return
			}
		}
		t.Hash = state.Sum(nil)

	default:
		panic(fmt.Errorf("unexpected token: %+v", token))

	}

	return
}
