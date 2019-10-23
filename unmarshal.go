package sb

import (
	"encoding"
	"reflect"
)

func Unmarshal(stream Stream, target any) error {
	return UnmarshalValue(stream, reflect.ValueOf(target))
}

type Detokenizer interface {
	DetokenizeSB(stream Stream) error
}

var (
	detokenizerType       = reflect.TypeOf((*Detokenizer)(nil)).Elem()
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func UnmarshalValue(stream Stream, ptr reflect.Value) error {

	if ptr.IsValid() {
		if t := ptr.Type(); t.Implements(detokenizerType) {
			return ptr.Interface().(Detokenizer).DetokenizeSB(stream)
		} else if t.Implements(binaryUnmarshalerType) {
			p, err := stream.Next()
			if err != nil {
				return err
			}
			if p == nil {
				return UnmarshalError{ExpectingString}
			}
			token := *p
			if token.Kind != KindString {
				return UnmarshalError{ExpectingString}
			}
			if err = ptr.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(
				[]byte(token.Value.(string)),
			); err != nil {
				return err
			}
			return nil
		} else if t.Implements(textUnmarshalerType) {
			p, err := stream.Next()
			if err != nil {
				return err
			}
			if p == nil {
				return UnmarshalError{ExpectingString}
			}
			token := *p
			if token.Kind != KindString {
				return UnmarshalError{ExpectingString}
			}
			if err = ptr.Interface().(encoding.TextUnmarshaler).UnmarshalText(
				[]byte(token.Value.(string)),
			); err != nil {
				return err
			}
			return nil
		}
	}

	p, err := stream.Next()
	if err != nil {
		return err
	}
	if p == nil {
		return UnmarshalError{ExpectingValue}
	}
	token := *p

	switch token.Kind {
	case KindNil, KindArrayEnd, KindObjectEnd:
		return nil
	}

	valueType := ptr.Type().Elem()
	hasConcreteType := false
	if valueType.Kind() != reflect.Interface {
		hasConcreteType = true
		if ptr.IsNil() {
			valuePtr := reflect.New(valueType)
			ptr.Elem().Set(valuePtr.Elem())
		}
	}

	switch token.Kind {

	case KindBool:
		if hasConcreteType {
			ptr.Elem().SetBool(token.Value.(bool))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(bool)))
		}

	case KindInt:
		if hasConcreteType {
			ptr.Elem().SetInt(int64(token.Value.(int)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int)))
		}

	case KindInt8:
		if hasConcreteType {
			ptr.Elem().SetInt(int64(token.Value.(int8)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int8)))
		}

	case KindInt16:
		if hasConcreteType {
			ptr.Elem().SetInt(int64(token.Value.(int16)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int16)))
		}

	case KindInt32:
		if hasConcreteType {
			ptr.Elem().SetInt(int64(token.Value.(int32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int32)))
		}

	case KindInt64:
		if hasConcreteType {
			ptr.Elem().SetInt(token.Value.(int64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int64)))
		}

	case KindUint:
		if hasConcreteType {
			ptr.Elem().SetUint(uint64(token.Value.(uint)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint)))
		}

	case KindUint8:
		if hasConcreteType {
			ptr.Elem().SetUint(uint64(token.Value.(uint8)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
		}

	case KindUint16:
		if hasConcreteType {
			ptr.Elem().SetUint(uint64(token.Value.(uint16)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
		}

	case KindUint32:
		if hasConcreteType {
			ptr.Elem().SetUint(uint64(token.Value.(uint32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
		}

	case KindUint64:
		if hasConcreteType {
			ptr.Elem().SetUint(token.Value.(uint64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
		}

	case KindFloat32:
		if hasConcreteType {
			ptr.Elem().SetFloat(float64(token.Value.(float32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float32)))
		}

	case KindFloat64:
		if hasConcreteType {
			ptr.Elem().SetFloat(token.Value.(float64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float64)))
		}

	case KindString:
		if hasConcreteType {
			ptr.Elem().SetString(token.Value.(string))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(string)))
		}

	case KindArray:
		if hasConcreteType {

			if valueType.Kind() == reflect.Array {
				// array
				idx := 0
				for {
					p, err := stream.Peek()
					if err != nil {
						return err
					}
					if p == nil {
						return UnmarshalError{ExpectingValue}
					}
					if p.Kind == KindArrayEnd {
						stream.Next()
						break
					}
					err = UnmarshalValue(stream, ptr.Elem().Index(idx).Addr())
					if err != nil {
						return err
					}
					idx++
				}

			} else {
				// slice
				slice := ptr.Elem()
				for {
					p, err := stream.Peek()
					if err != nil {
						return err
					}
					if p == nil {
						return UnmarshalError{ExpectingValue}
					}
					if p.Kind == KindArrayEnd {
						stream.Next()
						break
					}
					elemPtr := reflect.New(valueType.Elem())
					err = UnmarshalValue(stream, elemPtr)
					if err != nil {
						return err
					}
					slice = reflect.Append(slice, elemPtr.Elem())
				}
				ptr.Elem().Set(slice)
			}

		} else {
			// generic slice
			var slice []any
			for {
				p, err := stream.Peek()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindArrayEnd {
					stream.Next()
					break
				}
				var elem any
				err = UnmarshalValue(stream, reflect.ValueOf(&elem))
				if err != nil {
					return err
				}
				slice = append(slice, elem)
			}
			ptr.Elem().Set(reflect.ValueOf(slice))
		}

	case KindObject:
		if hasConcreteType {
			for {
				p, err := stream.Peek()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindObjectEnd {
					stream.Next()
					break
				}

				// name
				var name string
				err = UnmarshalValue(stream, reflect.ValueOf(&name))
				if err != nil {
					return err
				}

				// value
				field, ok := valueType.FieldByName(name)
				if !ok || field.Anonymous {
					// skip next value
					var value any
					err = UnmarshalValue(stream, reflect.ValueOf(&value))
					if err != nil {
						return err
					}
					continue
				}
				err = UnmarshalValue(stream, ptr.Elem().FieldByIndex(field.Index).Addr())
				if err != nil {
					return err
				}
			}

		} else {
			// construct new type
			var values []any
			var fields []reflect.StructField
			for {
				p, err := stream.Peek()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindObjectEnd {
					stream.Next()
					break
				}

				// name
				var name string
				err = UnmarshalValue(stream, reflect.ValueOf(&name))
				if err != nil {
					return err
				}

				// value
				var value any
				err = UnmarshalValue(stream, reflect.ValueOf(&value))
				if err != nil {
					return err
				}
				values = append(values, value)
				fields = append(fields, reflect.StructField{
					Name: name,
					Type: reflect.TypeOf(value),
				})
			}

			structType := reflect.StructOf(fields)
			structPtr := reflect.New(structType)
			for i, value := range values {
				structPtr.Elem().Field(i).Set(reflect.ValueOf(value))
			}
			ptr.Elem().Set(structPtr.Elem())
		}

	}

	return nil
}
