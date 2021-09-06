package sb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"

	"github.com/reusee/e4"
)

func DecodeBuffer(r io.Reader, byteReader io.ByteReader, buf []byte, cont Proc) Proc {
	return decodeBuffer(r, byteReader, buf, false, cont)
}

func DecodeBufferForCompare(r io.Reader, byteReader io.ByteReader, buf []byte, cont Proc) Proc {
	return decodeBuffer(r, byteReader, buf, true, cont)
}

var initDecodeStep = 8

var MaxDecodeStringLength uint64 = 4 * 1024 * 1024 * 1024

func decodeBuffer(r io.Reader, byteReader io.ByteReader, buf []byte, forCompare bool, cont Proc) Proc {
	var proc Proc
	var offset int64
	proc = Proc(func() (token *Token, next Proc, err error) {
		var kind Kind
		if byteReader != nil {
			if b, err := byteReader.ReadByte(); errors.Is(err, io.EOF) {
				return nil, cont, nil
			} else if err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset++
				kind = Kind(b)
			}
		} else {
			if _, err := io.ReadFull(r, buf[:1]); errors.Is(err, io.EOF) {
				return nil, cont, nil
			} else if err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
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
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				}
				offset++
				if b > 0 {
					value = true
				} else {
					value = false
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
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
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 8
			}
			value = int(binary.LittleEndian.Uint64(buf[:8]))

		case KindInt8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset++
					value = int8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset += 1
				}
				value = int8(buf[0])
			}

		case KindInt16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 2
			}
			value = int16(binary.LittleEndian.Uint16(buf[:2]))

		case KindInt32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 4
			}
			value = int32(binary.LittleEndian.Uint32(buf[:4]))

		case KindInt64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 8
			}
			value = int64(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 8
			}
			value = uint(binary.LittleEndian.Uint64(buf[:8]))

		case KindUint8:
			if byteReader != nil {
				if b, err := byteReader.ReadByte(); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset++
					value = uint8(b)
				}
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset += 1
				}
				value = uint8(buf[0])
			}

		case KindUint16:
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 2
			}
			value = binary.LittleEndian.Uint16(buf[:2])

		case KindUint32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 4
			}
			value = binary.LittleEndian.Uint32(buf[:4])

		case KindUint64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 8
			}
			value = binary.LittleEndian.Uint64(buf[:8])

		case KindFloat32:
			if _, err := io.ReadFull(r, buf[:4]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 4
			}
			value = math.Float32frombits(binary.LittleEndian.Uint32(buf[:4]))

		case KindFloat64:
			if _, err := io.ReadFull(r, buf[:8]); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += 8
			}
			value = math.Float64frombits(binary.LittleEndian.Uint64(buf[:8]))

		case KindString, KindTypeName, KindLiteral:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				}
				offset++
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
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
					return nil, nil, we(e4.With(Offset(offset)), e4.With(StringTooLong))(DecodeError)
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset += int64(l)
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				}
			}
			if length > MaxDecodeStringLength {
				return nil, nil, we(e4.With(Offset(offset)), e4.With(StringTooLong))(DecodeError)
			}

			if forCompare {
				length := int(length)
				step := initDecodeStep
				var segments func() (*Token, Proc, error)
				segments = func() (*Token, Proc, error) {
					if length == 0 {
						return &Token{
							Kind: KindStringEnd,
						}, proc, nil
					}
					l := step
					step *= 2
					if l > length {
						l = length
					}
					length -= l
					builder := new(strings.Builder)
					builder.Grow(l)
					buf, put := bytesPool32K.Get()
					defer put()
					if n, err := io.CopyBuffer(
						builder,
						io.LimitReader(r, int64(l)),
						*buf.(*[]byte),
					); err != nil || n != int64(l) {
						return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
					} else {
						offset += int64(length)
					}
					token := &Token{
						Kind:  kind,
						Value: builder.String(),
					}
					return token, segments, nil
				}
				return &Token{
					Kind: KindStringBegin,
				}, segments, nil
			}

			builder := new(strings.Builder)
			builder.Grow(int(length))
			buf, put := bytesPool32K.Get()
			defer put()
			if n, err := io.CopyBuffer(
				builder,
				io.LimitReader(r, int64(length)),
				*buf.(*[]byte),
			); err != nil || n != int64(length) {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
			} else {
				offset += int64(length)
			}
			value = builder.String()

		case KindBytes, KindRef:
			var length uint64
			var b byte
			if byteReader != nil {
				if b, err = byteReader.ReadByte(); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				}
				offset++
			} else {
				if _, err := io.ReadFull(r, buf[:1]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
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
					return nil, nil, we(e4.With(Offset(offset)), e4.With(BytesTooLong))(DecodeError)
				}
				if _, err := io.ReadFull(r, buf[:l]); err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				} else {
					offset += int64(l)
				}
				var err error
				length, err = binary.ReadUvarint(bytes.NewReader(buf[:l]))
				if err != nil {
					return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
				}
			}
			if length > MaxDecodeStringLength {
				return nil, nil, we(e4.With(Offset(offset)), e4.With(BytesTooLong))(DecodeError)
			}

			if forCompare {
				length := int(length)
				step := initDecodeStep
				var segments func() (*Token, Proc, error)
				segments = func() (*Token, Proc, error) {
					if length == 0 {
						return &Token{
							Kind: KindBytesEnd,
						}, proc, nil
					}
					l := step
					step *= 2
					if l > length {
						l = length
					}
					length -= l
					builder := new(bytes.Buffer)
					builder.Grow(l)
					buf, put := bytesPool32K.Get()
					defer put()
					if n, err := io.CopyBuffer(
						builder,
						io.LimitReader(r, int64(l)),
						*buf.(*[]byte),
					); err != nil || n != int64(l) {
						return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
					} else {
						offset += int64(length)
					}
					token := &Token{
						Kind:  kind,
						Value: builder.Bytes(),
					}
					return token, segments, nil
				}
				return &Token{
					Kind: KindBytesBegin,
				}, segments, nil
			}

			bs := make([]byte, length)
			if _, err := io.ReadFull(r, bs); err != nil {
				return nil, nil, we(e4.With(DecodeError), e4.With(Offset(offset)))(err)
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
			return nil, nil, we(e4.With(Offset(offset)), e4.With(BadTokenKind), e4.With(kind))(DecodeError)

		}

		return &Token{
			Kind:  kind,
			Value: value,
		}, proc, nil
	})

	return proc
}

func decode(r io.Reader, forCompare bool) Proc {
	var byteReader io.ByteReader
	if rd, ok := r.(io.ByteReader); ok {
		byteReader = rd
	}
	return decodeBuffer(r, byteReader, make([]byte, 8), forCompare, nil)
}

func Decode(r io.Reader) Proc {
	return decode(r, false)
}

func DecodeForCompare(r io.Reader) Proc {
	return decode(r, true)
}
