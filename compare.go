package sb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/reusee/e4"
)

func Compare(stream1, stream2 Stream) (int, error) {
	for {
		t1, err := stream1.Next()
		if err != nil {
			return 0, err
		}
		t2, err := stream2.Next()
		if err != nil {
			return 0, err
		}

		if t1 == nil && t2 == nil {
			break
		} else if t1 == nil && t2 != nil {
			return -1, nil
		} else if t1 != nil && t2 == nil {
			return 1, nil
		}

		if t1.Kind != t2.Kind {
			if t1.Kind < t2.Kind {
				return -1, nil
			} else {
				return 1, nil
			}
		}

		if v1, ok := t1.Value.([]byte); ok {
			if v2, ok := t2.Value.([]byte); ok {
				if res := bytes.Compare(v1, v2); res != 0 {
					return res, nil
				}
			}

		} else {
			if t1.Value != t2.Value {
				switch v1 := t1.Value.(type) {

				case bool:
					v2 := t2.Value.(bool)
					if !v1 && v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case int:
					v2 := t2.Value.(int)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case int8:
					v2 := t2.Value.(int8)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case int16:
					v2 := t2.Value.(int16)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case int32:
					v2 := t2.Value.(int32)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case int64:
					v2 := t2.Value.(int64)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case uint:
					v2 := t2.Value.(uint)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case uint8:
					v2 := t2.Value.(uint8)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case uint16:
					v2 := t2.Value.(uint16)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case uint32:
					v2 := t2.Value.(uint32)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case uint64:
					v2 := t2.Value.(uint64)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case float32:
					v2 := t2.Value.(float32)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case float64:
					v2 := t2.Value.(float64)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				case string:
					v2 := t2.Value.(string)
					if v1 < v2 {
						return -1, nil
					} else {
						return 1, nil
					}

				default: // NOCOVER
					panic(fmt.Errorf("bad type %#v %T", v1, v1))

				}
			}

		}

	}

	return 0, nil
}

func MustCompare(stream1, stream2 Stream) int {
	res, err := Compare(stream1, stream2)
	if err != nil { // NOCOVER
		panic(err)
	}
	return res
}

func CompareBytes(a, b []byte) (int, error) {

	var offsetA int64
	var offsetB int64
	readA := func(l int) (ret []byte, err error) {
		if len(a) < l {
			return nil, we(e4.With(Offset(offsetA)), e4.With(io.ErrUnexpectedEOF))(DecodeError)
		}
		ret = a[:l]
		a = a[l:]
		offsetA += int64(l)
		return
	}
	readB := func(l int) (ret []byte, err error) {
		if len(b) < l {
			return nil, we(e4.With(Offset(offsetB)), e4.With(io.ErrUnexpectedEOF))(DecodeError)
		}
		ret = b[:l]
		b = b[l:]
		offsetB += int64(l)
		return
	}

	for {

		if len(a) == 0 && len(b) == 0 {
			return 0, nil
		}
		if len(a) == 0 && len(b) != 0 {
			return -1, nil
		}
		if len(a) != 0 && len(b) == 0 {
			return 1, nil
		}

		bs, err := readA(1)
		if err != nil {
			return 0, err
		}
		kindA := Kind(bs[0])
		bs, err = readB(1)
		if err != nil {
			return 0, err
		}
		kindB := Kind(bs[0])
		if kindA < kindB {
			return -1, nil
		}
		if kindA > kindB {
			return 1, nil
		}

		switch kindA {

		case KindBool:
			bs, err = readA(1)
			if err != nil {
				return 0, err
			}
			a1 := bs[0] > 0
			bs, err = readB(1)
			if err != nil {
				return 0, err
			}
			b1 := bs[0] > 0
			if !a1 && b1 {
				return -1, nil
			} else if a1 && !b1 {
				return 1, nil
			}

		case KindInt, KindInt64, KindUint, KindUint64:
			bs, err = readA(8)
			if err != nil {
				return 0, err
			}
			a1 := binary.LittleEndian.Uint64(bs)
			bs, err = readB(8)
			if err != nil {
				return 0, err
			}
			b1 := binary.LittleEndian.Uint64(bs)
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindInt8, KindUint8:
			bs, err = readA(1)
			if err != nil {
				return 0, err
			}
			a1 := bs[0]
			bs, err = readB(1)
			if err != nil {
				return 0, err
			}
			b1 := bs[0]
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindInt16, KindUint16:
			bs, err = readA(2)
			if err != nil {
				return 0, err
			}
			a1 := binary.LittleEndian.Uint16(bs)
			bs, err = readB(2)
			if err != nil {
				return 0, err
			}
			b1 := binary.LittleEndian.Uint16(bs)
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindInt32, KindUint32:
			bs, err = readA(4)
			if err != nil {
				return 0, err
			}
			a1 := binary.LittleEndian.Uint32(bs)
			bs, err = readB(4)
			if err != nil {
				return 0, err
			}
			b1 := binary.LittleEndian.Uint32(bs)
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindFloat32:
			bs, err := readA(4)
			if err != nil {
				return 0, err
			}
			a1 := math.Float32frombits(binary.LittleEndian.Uint32(bs))
			bs, err = readB(4)
			if err != nil {
				return 0, err
			}
			b1 := math.Float32frombits(binary.LittleEndian.Uint32(bs))
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindFloat64:
			bs, err := readA(8)
			if err != nil {
				return 0, err
			}
			a1 := math.Float64frombits(binary.LittleEndian.Uint64(bs))
			bs, err = readB(8)
			if err != nil {
				return 0, err
			}
			b1 := math.Float64frombits(binary.LittleEndian.Uint64(bs))
			if a1 < b1 {
				return -1, nil
			} else if a1 > b1 {
				return 1, nil
			}

		case KindString, KindBytes, KindTypeName:
			var l1 int
			bs, err := readA(1)
			if err != nil {
				return 0, err
			}
			x := bs[0]
			if x < 128 {
				l1 = int(x)
			} else {
				l := ^x
				if l > 8 {
					return 0, we(e4.With(Offset(offsetA)), e4.With(StringTooLong))(DecodeError)
				}
				bs, err = readA(int(l))
				if err != nil {
					return 0, err
				}
				n, _ := binary.Uvarint(bs)
				if n == 0 {
					return 0, we(e4.With(Offset(offsetA)), e4.With(BadStringLength))(DecodeError)
				}
				l1 = int(n)
			}
			a1, err := readA(l1)
			if err != nil {
				return 0, err
			}
			var l2 int
			bs, err = readB(1)
			if err != nil {
				return 0, err
			}
			x = bs[0]
			if x < 128 {
				l2 = int(x)
			} else {
				l := ^x
				if l > 8 {
					return 0, we(e4.With(Offset(offsetB)), e4.With(StringTooLong))(DecodeError)
				}
				bs, err = readB(int(l))
				if err != nil {
					return 0, err
				}
				n, _ := binary.Uvarint(bs)
				if n == 0 {
					return 0, we(e4.With(Offset(offsetB)), e4.With(BadStringLength))(DecodeError)
				}
				l2 = int(n)
			}
			b1, err := readB(l2)
			if err != nil {
				return 0, err
			}
			if res := bytes.Compare(a1, b1); res != 0 {
				return res, nil
			}

		case KindMin,
			KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd,
			KindNil, KindNaN,
			KindArray, KindObject, KindMap, KindTuple,
			KindMax:

		default:
			return 0, we(e4.With(Offset(offsetA)), e4.With(BadTokenKind), e4.With(kindA))(DecodeError)

		}

	}

}

func MustCompareBytes(a, b []byte) int {
	res, err := CompareBytes(a, b)
	if err != nil {
		panic(err)
	}
	return res
}
