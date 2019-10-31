package sb

import (
	"fmt"
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

			default:
				panic(fmt.Errorf("bad type %#v %T", v1, v1))

			}
		}

	}

	return 0, nil
}

func MustCompare(stream1, stream2 Stream) int {
	res, err := Compare(stream1, stream2)
	if err != nil {
		panic(err)
	}
	return res
}
