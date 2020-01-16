package sb

import (
	"bytes"
	"fmt"
)

func Compare(stream1, stream2 Stream) (int, error) {
	for {
	read1:
		t1, err := stream1.Next()
		if err != nil {
			return 0, err
		}
		if t1 != nil && t1.Kind == KindPostHash {
			goto read1
		}
	read2:
		t2, err := stream2.Next()
		if err != nil {
			return 0, err
		}
		if t2 != nil && t2.Kind == KindPostHash {
			goto read2
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
