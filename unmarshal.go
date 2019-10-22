package sb

import "reflect"

func Unmarshal(tokenizer Tokenizer, target any) error {
	_, err := UnmarshalValue(tokenizer, reflect.ValueOf(target))
	return err
}

func UnmarshalValue(tokenizer Tokenizer, ptr reflect.Value) (token Token, err error) {
	p := tokenizer.Next()
	if p == nil {
		return
	}
	token = *p
	switch token.Kind {

	case KindBool:
		ptr.Elem().Set(reflect.ValueOf(token.Value.(bool)))

	case KindInt:
		ptr.Elem().Set(reflect.ValueOf(token.Value.(int64)))

	case KindUint:
		ptr.Elem().Set(reflect.ValueOf(token.Value.(uint64)))

	case KindFloat:
		ptr.Elem().Set(reflect.ValueOf(token.Value.(float64)))

	case KindString:
		ptr.Elem().Set(reflect.ValueOf(token.Value.(string)))

	case KindNil, KindArrayEnd, KindObjectEnd:
		return

	case KindArray:
		var arrayPtr reflect.Value
		switch ptr.Type().Elem().Kind() {
		case reflect.Slice:
			arrayPtr = ptr
		default:
			array := []any{}
			arrayPtr = reflect.ValueOf(&array)
		}
		elemType := arrayPtr.Type().Elem().Elem()
		for {
			elemPtr := reflect.New(elemType)
			var subToken Token
			subToken, err = UnmarshalValue(tokenizer, elemPtr)
			if err != nil {
				return
			}
			if subToken.Kind == KindArrayEnd {
				break
			}
			arrayPtr.Elem().Set(
				reflect.Append(
					arrayPtr.Elem(),
					elemPtr.Elem(),
				),
			)
		}
		ptr.Elem().Set(arrayPtr.Elem())

	case KindObject:
		var values []any
		var fields []reflect.StructField
		newType := false
		if ptr.Type().Elem().Kind() != reflect.Struct {
			newType = true
		}

		for {

			// name
			var name string
			var subToken Token
			subToken, err = UnmarshalValue(tokenizer, reflect.ValueOf(&name))
			if err != nil {
				return
			}
			if subToken.Kind == KindObjectEnd {
				break
			}

			// value
			var value any
			_, err = UnmarshalValue(tokenizer, reflect.ValueOf(&value))
			if err != nil {
				return
			}

			// set or save
			if newType {
				values = append(values, value)
				fields = append(fields, reflect.StructField{
					Name: name,
					Type: reflect.TypeOf(value),
				})
			} else {
				field, ok := ptr.Type().Elem().FieldByName(name)
				if !ok {
					continue
				}
				ptr.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(value))
			}

		}

		// create new type
		structType := reflect.StructOf(fields)
		structPtr := reflect.New(structType)
		for i, value := range values {
			structPtr.Elem().Field(i).Set(reflect.ValueOf(value))
		}
		ptr.Elem().Set(structPtr.Elem())

	}

	return
}
