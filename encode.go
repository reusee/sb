package sb

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Encode(w io.Writer, tokenizer Tokenizer) error {
	buf := make([]byte, 8)
	for token := tokenizer.Next(); token != nil; token = tokenizer.Next() {
		if err := binary.Write(w, binary.LittleEndian, token.Kind); err != nil {
			return err
		}
		if token.Value != nil {
			switch value := token.Value.(type) {

			case bool:
				if value {
					if _, err := w.Write([]byte{1}); err != nil {
						return err
					}
				} else {
					if _, err := w.Write([]byte{0}); err != nil {
						return err
					}
				}

			case int:
				if err := binary.Write(w, binary.LittleEndian, int64(value)); err != nil {
					return err
				}

			case uint:
				if err := binary.Write(w, binary.LittleEndian, uint64(value)); err != nil {
					return err
				}

			case int8, int16, int32, int64,
				uint8, uint16, uint32, uint64,
				float32, float64:
				if err := binary.Write(w, binary.LittleEndian, value); err != nil {
					return err
				}

			case string:
				l := uint64(len(value))
				if l < 128 {
					if _, err := w.Write([]byte{byte(l)}); err != nil {
						return err
					}
				} else {
					n := binary.PutUvarint(buf, l)
					if _, err := w.Write([]byte{byte(^n)}); err != nil {
						return err
					}
					if _, err := w.Write(buf[:n]); err != nil {
						return err
					}
				}
				if _, err := w.Write([]byte(value)); err != nil {
					return err
				}

			default:
				panic(DecodeError(fmt.Errorf("bad type %#v %T", value, value)))

			}
		}
	}

	return nil
}
