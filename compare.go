package sb

func Compare(a, b any) int {
	tokenizer1 := NewTokenizer(a)
	tokenizer2 := NewTokenizer(b)
	for {
		t1 := tokenizer1.Next()
		t2 := tokenizer2.Next()

		if t1 == nil && t2 == nil {
			break
		} else if t1 == nil && t2 != nil {
			return -1
		} else if t1 != nil && t1 == nil {
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
			case int64:
				v2 := t2.Value.(int64)
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
			}
		}

	}

	return 0
}
