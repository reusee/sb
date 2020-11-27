package sb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"
	"sync"
)

var copyBufferPool = sync.Pool{
	New: func() any {
		bs := make([]byte, 32*1024)
		return &bs
	},
}

func DecodeBuffer(r io.Reader, byteReader io.ByteReader, buf []byte, cont Proc) Proc {
	var proc Proc
	var offset int64
	proc = Proc(func() (token *Token, next Proc, err error) {
		defer func() {
			if next == nil {
				decodeBufPool.Put(&buf)
			}
		}()
		var kind Kind
		if byteReader != nil {
			if b, err := byteReader.ReadByte(); errors.Is(err, io.EOF) {
				return nil, nil, nil
			} else if err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset++
				kind = Kind(b)
			}
		} else {
			if _, err := io.ReadFull(r, buf[:1]); errors.Is(err, io.EOF) {
				return nil, nil, nil
			} else if err != nil {
				return nil, nil, NewDecodeError(offset, err)
			}
			offset += 1
			kind = Kind(buf[0])
		}

		var value any
		switch kind {

		case KindBool:
			if byteReader != nil {
				b, err := byteReader.ReadByte()
				if err != nil {
					return nil, nil, NewDecodeError(offset, err)
				}
				offset++
				if b > 0 {
					value = true
				} else {
					value = false
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += 1
				}
				if buf[0] > 0 {
					value = true
				} else {
					value = false
				}
			}

		case KindInt:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 8
			}
			value = int(binary.LittleEndian.Uint64(buf[:8]))

		case KindInt8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset++
					value = int8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += 1
				}
				value = int8(buf[0])
			}

		case KindInt16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 2
			}
			value = int16(binary.LittleEndian.Uint16(buf[:2]))

		case KindInt32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 4
			}
			value = int32(binary.LittleEndian.Uint32(buf[:4]))

		case KindInt64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 8
			}
			value = int64(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 8
			}
			value = uint(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset++
					value = uint8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += 1
				}
				value = uint8(buf[0])
			}

		case KindUint16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 2
			}
			value = binary.LittleEndian.Uint16(buf[:2])

		case KindUint32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 4
			}
			value = binary.LittleEndian.Uint32(buf[:4])

		case KindUint64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 8
			}
			value = binary.LittleEndian.Uint64(buf[:8])

		case KindFloat32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 4
			}
			value = math.Float32frombits(binary.LittleEndian.Uint32(buf[:4]))

		case KindFloat64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += 8
			}
			value = math.Float64frombits(binary.LittleEndian.Uint64(buf[:8]))

		case KindString:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				}
				offset++
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += 1
				}
				b = buf[0]
			}
			if b < 128 {
				length = uint64(b)
			} else {
				l := ^b
				if l > 8 {
					return nil, nil, NewDecodeError(offset, StringTooLong)
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += int64(l)
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, NewDecodeError(offset, err)
				}
			}
			if length > 128*1024*1024 {
				return nil, nil, NewDecodeError(offset, StringTooLong)
			}
			builder := new(strings.Builder)
			builder.Grow(int(length))
			buf := copyBufferPool.Get().(*[]byte)
			defer copyBufferPool.Put(buf)
			if n, err := io.CopyBuffer(
				builder,
				io.LimitReader(r, int64(length)),
				*buf,
			); err != nil || n != int64(length) {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += int64(length)
			}
			value = builder.String()

		case KindBytes:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				}
				offset++
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += 1
				}
				b = buf[0]
			}
			if b < 128 {
				length = uint64(b)
			} else {
				l := ^b
				if l > 8 {
					return nil, nil, NewDecodeError(offset, BytesTooLong)
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, NewDecodeError(offset, err)
				} else {
					offset += int64(l)
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, NewDecodeError(offset, err)
				}
			}
			if length > 128*1024*1024 {
				return nil, nil, NewDecodeError(offset, BytesTooLong)
			}
			bs := make([]byte, length)
			if _, err := io.ReadFull(r, bs); err != nil {
				return nil, nil, NewDecodeError(offset, err)
			} else {
				offset += int64(length)
			}
			value = bs

		case KindMin,
			KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd,
			KindNil, KindNaN,
			KindArray, KindObject, KindMap, KindTuple,
			KindMax:

		default:
			return nil, nil, NewDecodeError(offset, BadTokenKind, kind)

		}

		return &Token{
			Kind:  kind,
			Value: value,
		}, proc, nil
	})

	return proc
}

var decodeBufPool = sync.Pool{
	New: func() any {
		bs := make([]byte, 8)
		return &bs
	},
}

func Decode(r io.Reader) *Proc {
	var byteReader io.ByteReader
	if rd, ok := r.(io.ByteReader); ok {
		byteReader = rd
	}
	buf := decodeBufPool.Get().(*[]byte)
	proc := DecodeBuffer(r, byteReader, *buf, nil)
	return &proc
}
