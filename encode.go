package sb

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func Encode(w io.Writer) Sink {
	return EncodeBuffer(
		w,
		make([]byte, 8),
		nil,
	)
}

func EncodeBuffer(w io.Writer, buf []byte, cont Sink) Sink {
	var byteWriter io.ByteWriter
	if bw, ok := w.(io.ByteWriter); ok {
		byteWriter = bw
	}

	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token.Invalid() {
			return cont, nil
		}
		if byteWriter != nil {
			if err := byteWriter.WriteByte(byte(token.Kind)); err != nil { // NOCOVER
				return nil, err
			}
		} else {
			if _, err := w.Write([]byte{byte(token.Kind)}); err != nil {
				return nil, err
			}
		}
		if token.Value != nil {
			switch value := token.Value.(type) {

			case bool:
				if value {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(1); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{1}); err != nil { // NOCOVER
							return nil, err
						}
					}
				} else {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(0); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{0}); err != nil { // NOCOVER
							return nil, err
						}
					}
				}

			case int:
				binary.LittleEndian.PutUint64(buf, uint64(int64(value)))
				if _, err := w.Write(buf); err != nil {
					return nil, err
				}

			case int8:
				if _, err := w.Write([]byte{uint8(value)}); err != nil {
					return nil, err
				}

			case int16:
				binary.LittleEndian.PutUint16(buf, uint16(value))
				if _, err := w.Write((buf)[:2]); err != nil {
					return nil, err
				}

			case int32:
				binary.LittleEndian.PutUint32(buf, uint32(value))
				if _, err := w.Write((buf)[:4]); err != nil {
					return nil, err
				}

			case int64:
				binary.LittleEndian.PutUint64(buf, uint64(value))
				if _, err := w.Write((buf)); err != nil {
					return nil, err
				}

			case uint:
				binary.LittleEndian.PutUint64(buf, uint64(value))
				if _, err := w.Write(buf); err != nil {
					return nil, err
				}

			case uintptr:
				binary.LittleEndian.PutUint64(buf, uint64(value))
				if _, err := w.Write(buf); err != nil {
					return nil, err
				}

			case uint8:
				if _, err := w.Write([]byte{value}); err != nil {
					return nil, err
				}

			case uint16:
				binary.LittleEndian.PutUint16(buf, value)
				if _, err := w.Write((buf)[:2]); err != nil {
					return nil, err
				}

			case uint32:
				binary.LittleEndian.PutUint32(buf, value)
				if _, err := w.Write((buf)[:4]); err != nil {
					return nil, err
				}

			case uint64:
				binary.LittleEndian.PutUint64(buf, value)
				if _, err := w.Write((buf)); err != nil {
					return nil, err
				}

			case float32:
				binary.LittleEndian.PutUint32(buf, math.Float32bits(value))
				if _, err := w.Write((buf)[:4]); err != nil {
					return nil, err
				}

			case float64:
				binary.LittleEndian.PutUint64(buf, math.Float64bits(value))
				if _, err := w.Write((buf)); err != nil {
					return nil, err
				}

			case string:
				l := uint64(len(value))
				if l < 128 {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(byte(l)); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{byte(l)}); err != nil { // NOCOVER
							return nil, err
						}
					}
				} else {
					n := binary.PutUvarint(buf, l)
					if byteWriter != nil {
						if err := byteWriter.WriteByte(byte(^n)); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{byte(^n)}); err != nil { // NOCOVER
							return nil, err
						}
					}
					if _, err := w.Write(buf[:n]); err != nil { // NOCOVER
						return nil, err
					}
				}
				if _, err := w.Write([]byte(value)); err != nil { // NOCOVER
					return nil, err
				}

			case []byte:
				l := uint64(len(value))
				if l < 128 {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(byte(l)); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{byte(l)}); err != nil { // NOCOVER
							return nil, err
						}
					}
				} else {
					n := binary.PutUvarint(buf, l)
					if byteWriter != nil {
						if err := byteWriter.WriteByte(byte(^n)); err != nil { // NOCOVER
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{byte(^n)}); err != nil { // NOCOVER
							return nil, err
						}
					}
					if _, err := w.Write(buf[:n]); err != nil { // NOCOVER
						return nil, err
					}
				}
				if _, err := w.Write(value); err != nil { // NOCOVER
					return nil, err
				}

			default: // NOCOVER
				panic(fmt.Errorf("bad type %#v %T", value, value))

			}
		}

		return sink, nil
	}

	return sink
}
