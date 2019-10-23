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
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetBool(token.Value.(bool))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(bool)))
		}

	case KindInt:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetInt(int64(token.Value.(int)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int)))
		}

	case KindInt8:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetInt(int64(token.Value.(int8)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int8)))
		}

	case KindInt16:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetInt(int64(token.Value.(int16)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int16)))
		}

	case KindInt32:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetInt(int64(token.Value.(int32)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int32)))
		}

	case KindInt64:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetInt(int64(token.Value.(int64)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int64)))
		}

	case KindUint:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetUint(uint64(token.Value.(uint)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint)))
		}

	case KindUint8:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetUint(uint64(token.Value.(uint8)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
		}

	case KindUint16:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetUint(uint64(token.Value.(uint16)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
		}

	case KindUint32:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetUint(uint64(token.Value.(uint32)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
		}

	case KindUint64:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetUint(uint64(token.Value.(uint64)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
		}

	case KindFloat32:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetFloat(float64(token.Value.(float32)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float32)))
		}

	case KindFloat64:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetFloat(float64(token.Value.(float32)))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float64)))
		}

	case KindString:
		if ptr.Type().Elem().Name() != "" {
			// defined type
			v := reflect.New(ptr.Type().Elem()).Elem()
			v.SetString(token.Value.(string))
			ptr.Elem().Set(v)
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(string)))
		}

	case KindNil, KindArrayEnd, KindObjectEnd:
		return

	case KindArray:
		var arrayPtr reflect.Value
		isArray := false
		switch ptr.Type().Elem().Kind() {
		case reflect.Slice:
			arrayPtr = ptr
		case reflect.Array:
			sliceType := reflect.SliceOf(ptr.Type().Elem().Elem())
			arrayPtr = reflect.New(sliceType)
			isArray = true
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
		if isArray {
			aPtr := reflect.New(ptr.Type().Elem())
			reflect.Copy(aPtr.Elem().Slice(0, aPtr.Elem().Len()), arrayPtr.Elem())
			ptr.Elem().Set(aPtr.Elem())
		} else {
			ptr.Elem().Set(arrayPtr.Elem())
		}

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

			if newType {
				var value any
				_, err = UnmarshalValue(tokenizer, reflect.ValueOf(&value))
				if err != nil {
					return
				}
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
				valuePtr := reflect.New(field.Type)
				_, err = UnmarshalValue(tokenizer, valuePtr)
				if err != nil {
					return
				}
				ptr.Elem().FieldByIndex(field.Index).Set(valuePtr.Elem())
			}

		}

		if newType {
			structType := reflect.StructOf(fields)
			structPtr := reflect.New(structType)
			for i, value := range values {
				structPtr.Elem().Field(i).Set(reflect.ValueOf(value))
			}
			ptr.Elem().Set(structPtr.Elem())
		}

	}

	return
}