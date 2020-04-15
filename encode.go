package sb

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Encode(w io.Writer) Sink {
	buf := make([]byte, 8)
	return EncodeBuffer(w, buf, nil)
}

func EncodeBuffer(w io.Writer, buf []byte, cont Sink) Sink {
	var byteWriter io.ByteWriter
	if bw, ok := w.(io.ByteWriter); ok {
		byteWriter = bw
	}

	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return cont, nil
		}
		if byteWriter != nil {
			if err := byteWriter.WriteByte(byte(token.Kind)); err != nil {
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
						if err := byteWriter.WriteByte(1); err != nil {
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{1}); err != nil { // NOCOVER
							return nil, err
						}
					}
				} else {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(0); err != nil {
							return nil, err
						}
					} else {
						if _, err := w.Write([]byte{0}); err != nil { // NOCOVER
							return nil, err
						}
					}
				}

			case int:
				if err := binary.Write(w, binary.LittleEndian, int64(value)); err != nil { // NOCOVER
					return nil, err
				}

			case uint:
				if err := binary.Write(w, binary.LittleEndian, uint64(value)); err != nil { // NOCOVER
					return nil, err
				}

			case int8, int16, int32, int64,
				uint8, uint16, uint32, uint64,
				float32, float64:
				if err := binary.Write(w, binary.LittleEndian, value); err != nil { // NOCOVER
					return nil, err
				}

			case string:
				l := uint64(len(value))
				if l < 128 {
					if byteWriter != nil {
						if err := byteWriter.WriteByte(byte(l)); err != nil {
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
						if err := byteWriter.WriteByte(byte(^n)); err != nil {
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
						if err := byteWriter.WriteByte(byte(l)); err != nil {
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
						if err := byteWriter.WriteByte(byte(^n)); err != nil {
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
