package sb

import (
	"encoding/binary"
	"fmt"
)

func EncodedLen(ret *int, cont Sink) Sink {
	buf := make([]byte, 8)
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return cont, nil
		}
		*ret++
		if token.Value != nil {
			switch value := token.Value.(type) {

			case bool, int8, uint8:
				*ret += 1

			case int16, uint16:
				*ret += 2

			case int32, uint32, float32:
				*ret += 4

			case int, uint, int64, uint64, float64:
				*ret += 8

			case string:
				l := uint64(len(value))
				if l < 128 {
					*ret++
				} else {
					n := binary.PutUvarint(buf, l)
					*ret += 1 + n
				}
				*ret += len(value)

			case []byte:
				l := uint64(len(value))
				if l < 128 {
					*ret++
				} else {
					n := binary.PutUvarint(buf, l)
					*ret += 1 + n
				}
				*ret += len(value)

			default: // NOCOVER
				panic(fmt.Errorf("bad type %#v %T", value, value))

			}
		}

		return sink, nil
	}

	return sink
}
