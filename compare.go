package sb

import "fmt"

func Compare(tokenizer1, tokenizer2 Tokenizer) int {
	for {
		t1 := tokenizer1.Next()
		t2 := tokenizer2.Next()

		if t1 == nil && t2 == nil {
			break
		} else if t1 == nil && t2 != nil {
			return -1
		} else if t1 != nil && t2 == nil {
			return 1
		}

		if t1.Kind != t2.Kind {
			if t1.Kind < t2.Kind {
				return -1
			} else {
				return 1
			}
		}

		if t1.Value != t2.Value {
			switch v1 := t1.Value.(type) {

			case bool:
				v2 := t2.Value.(bool)
				if !v1 && v2 {
					return -1
				} else {
					return 1
				}

			case int:
				v2 := t2.Value.(int)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case int8:
				v2 := t2.Value.(int8)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case int16:
				v2 := t2.Value.(int16)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case int32:
				v2 := t2.Value.(int32)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case int64:
				v2 := t2.Value.(int64)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case uint:
				v2 := t2.Value.(uint)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case uint8:
				v2 := t2.Value.(uint8)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case uint16:
				v2 := t2.Value.(uint16)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case uint32:
				v2 := t2.Value.(uint32)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case uint64:
				v2 := t2.Value.(uint64)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case float32:
				v2 := t2.Value.(float32)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case float64:
				v2 := t2.Value.(float64)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			case string:
				v2 := t2.Value.(string)
				if v1 < v2 {
					return -1
				} else {
					return 1
				}

			default:
				panic(DecodeError(fmt.Errorf("bad type %#v %T", v1, v1)))

			}
		}

	}

	return 0
}
