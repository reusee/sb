package sb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

func Decode(r io.Reader) *Proc {
	var byteReader io.ByteReader
	if rd, ok := r.(io.ByteReader); ok {
		byteReader = rd
	}
	var proc Proc
	buf := make([]byte, 8)
	proc = Proc(func() (*Token, Proc, error) {
		var kind Kind
		if byteReader != nil {
			if b, err := byteReader.ReadByte(); errors.Is(err, io.EOF) {
				return nil, nil, nil
			} else if err != nil {
				return nil, nil, err
			} else {
				kind = Kind(b)
			}
		} else {
			if _, err := io.ReadFull(r, buf[:1]); errors.Is(err, io.EOF) {
				return nil, nil, nil
			} else if err != nil {
				return nil, nil, err
			}
			kind = Kind(buf[0])
		}

		var value any
		var err error
		switch kind {

		case KindBool:
			if byteReader != nil {
				b, err := byteReader.ReadByte()
				if err != nil {
					return nil, nil, err
				}
				if b > 0 {
					value = true
				} else {
					value = false
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, err
				}
				if buf[0] > 0 {
					value = true
				} else {
					value = false
				}
			}

		case KindInt:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, err
			}
			value = int(binary.LittleEndian.Uint64(buf[:8]))

		case KindInt8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, err
				} else {
					value = int8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, err
				}
				value = int8(buf[0])
			}

		case KindInt16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, err
			}
			value = int16(binary.LittleEndian.Uint16(buf[:2]))

		case KindInt32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, err
			}
			value = int32(binary.LittleEndian.Uint32(buf[:4]))

		case KindInt64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, err
			}
			value = int64(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, err
			}
			value = uint(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, err
				} else {
					value = uint8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, err
				}
				value = uint8(buf[0])
			}

		case KindUint16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, err
			}
			value = binary.LittleEndian.Uint16(buf[:2])

		case KindUint32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, err
			}
			value = binary.LittleEndian.Uint32(buf[:4])

		case KindUint64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, err
			}
			value = binary.LittleEndian.Uint64(buf[:8])

		case KindFloat32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, err
			}
			value = math.Float32frombits(binary.LittleEndian.Uint32(buf[:4]))

		case KindFloat64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, err
			}
			value = math.Float64frombits(binary.LittleEndian.Uint64(buf[:8]))

		case KindString:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, err
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, err
				}
				b = buf[0]
			}
			if b < 128 {
				length = uint64(b)
			} else {
				l := ^b
				if l > 8 {
					return nil, nil, DecodeError{StringTooLong}
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, err
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, err
				}
			}
			if length > 128*1024*1024 {
				return nil, nil, DecodeError{StringTooLong}
			}
			var bs []byte
			if int(length) <= len(buf) {
				bs = buf[:length]
			} else {
				bs = make([]byte, length)
			}
			if _, err := io.ReadFull(r, bs); err != nil {
				return nil, nil, err
			}
			value = string(bs)

		case KindBytes:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, err
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, err
				}
				b = buf[0]
			}
			if b < 128 {
				length = uint64(b)
			} else {
				l := ^b
				if l > 8 {
					return nil, nil, DecodeError{StringTooLong}
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, err
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, err
				}
			}
			if length > 128*1024*1024 {
				return nil, nil, DecodeError{StringTooLong}
			}
			bs := make([]byte, length)
			if _, err := io.ReadFull(r, bs); err != nil {
				return nil, nil, err
			}
			value = bs

		case KindMin,
			KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd,
			KindNil, KindNaN,
			KindArray, KindObject, KindMap, KindTuple,
			KindMax:

		default:
			return nil, nil, fmt.Errorf("%w: %d", BadTokenKind, kind)

		}

		return &Token{
			Kind:  kind,
			Value: value,
		}, proc, nil
	})

	return &proc
}
