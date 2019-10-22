package sb

import (
	"encoding/binary"
	"io"
)

func Encode(tokenizer Tokenizer, w io.Writer) error {
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

			case int64, uint64, float64:
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
				panic("bad value type")

			}
		}
	}

	return nil
}
