package sb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type Decoder struct {
	r   io.Reader
	buf [decoderBufLen]byte
}

const (
	decoderBufLen = 8
)

func Decode(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

var _ Stream = new(Decoder)

func (d *Decoder) Next() (token *Token, err error) {

	if _, err := io.ReadFull(d.r, d.buf[:1]); errors.Is(err, io.EOF) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	kind := Kind(d.buf[0])

	var value any
	switch kind {

	case KindBool:
		if _, err := io.ReadFull(d.r, d.buf[:1]); err != nil {
			return nil, err
		}
		if d.buf[0] > 0 {
			value = true
		} else {
			value = false
		}

	case KindInt:
		if _, err := io.ReadFull(d.r, d.buf[:8]); err != nil {
			return nil, err
		}
		value = int(binary.LittleEndian.Uint64(d.buf[:8]))

	case KindInt8:
		if _, err := io.ReadFull(d.r, d.buf[:1]); err != nil {
			return nil, err
		}
		value = int8(d.buf[0])

	case KindInt16:
		if _, err := io.ReadFull(d.r, d.buf[:2]); err != nil {
			return nil, err
		}
		value = int16(binary.LittleEndian.Uint16(d.buf[:2]))

	case KindInt32:
		if _, err := io.ReadFull(d.r, d.buf[:4]); err != nil {
			return nil, err
		}
		value = int32(binary.LittleEndian.Uint32(d.buf[:4]))

	case KindInt64:
		if _, err := io.ReadFull(d.r, d.buf[:8]); err != nil {
			return nil, err
		}
		value = int64(binary.LittleEndian.Uint64(d.buf[:8]))

	case KindUint:
		if _, err := io.ReadFull(d.r, d.buf[:8]); err != nil {
			return nil, err
		}
		value = uint(binary.LittleEndian.Uint64(d.buf[:8]))

	case KindUint8:
		if _, err := io.ReadFull(d.r, d.buf[:1]); err != nil {
			return nil, err
		}
		value = uint8(d.buf[0])

	case KindUint16:
		if _, err := io.ReadFull(d.r, d.buf[:2]); err != nil {
			return nil, err
		}
		value = binary.LittleEndian.Uint16(d.buf[:2])

	case KindUint32:
		if _, err := io.ReadFull(d.r, d.buf[:4]); err != nil {
			return nil, err
		}
		value = binary.LittleEndian.Uint32(d.buf[:4])

	case KindUint64:
		if _, err := io.ReadFull(d.r, d.buf[:8]); err != nil {
			return nil, err
		}
		value = binary.LittleEndian.Uint64(d.buf[:8])

	case KindFloat32:
		if _, err := io.ReadFull(d.r, d.buf[:4]); err != nil {
			return nil, err
		}
		value = math.Float32frombits(binary.LittleEndian.Uint32(d.buf[:4]))

	case KindFloat64:
		if _, err := io.ReadFull(d.r, d.buf[:8]); err != nil {
			return nil, err
		}
		value = math.Float64frombits(binary.LittleEndian.Uint64(d.buf[:8]))

	case KindString:
		var length uint64
		if _, err := io.ReadFull(d.r, d.buf[:1]); err != nil {
			return nil, err
		}
		if d.buf[0] < 128 {
			length = uint64(d.buf[0])
		} else {
			l := ^d.buf[0]
			if l > 8 {
				return nil, DecodeError{StringTooLong}
			}
			if _, err := io.ReadFull(d.r, d.buf[:l]); err != nil {
				return nil, err
			}
			var err error
			length, err = binary.ReadUvarint(bytes.NewReader(d.buf[:l]))
			if err != nil {
				return nil, err
			}
		}
		if length > 128*1024*1024 {
			return nil, DecodeError{StringTooLong}
		}
		var bs []byte
		if length <= decoderBufLen {
			bs = d.buf[:length]
		} else {
			bs = make([]byte, length)
		}
		if _, err := io.ReadFull(d.r, bs); err != nil {
			return nil, err
		}
		value = string(bs)

	case KindBytes, KindPostTag:
		var length uint64
		if _, err := io.ReadFull(d.r, d.buf[:1]); err != nil {
			return nil, err
		}
		if d.buf[0] < 128 {
			length = uint64(d.buf[0])
		} else {
			l := ^d.buf[0]
			if l > 8 {
				return nil, DecodeError{StringTooLong}
			}
			if _, err := io.ReadFull(d.r, d.buf[:l]); err != nil {
				return nil, err
			}
			var err error
			length, err = binary.ReadUvarint(bytes.NewReader(d.buf[:l]))
			if err != nil {
				return nil, err
			}
		}
		if length > 128*1024*1024 {
			return nil, DecodeError{StringTooLong}
		}
		bs := make([]byte, length)
		if _, err := io.ReadFull(d.r, bs); err != nil {
			return nil, err
		}
		value = bs

	case KindMin,
		KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd,
		KindNil, KindNaN,
		KindArray, KindObject, KindMap, KindTuple,
		KindMax:

	default:
		return nil, fmt.Errorf("%w: %d", BadTokenKind, kind)

	}

	return &Token{
		Kind:  kind,
		Value: value,
	}, nil
}
