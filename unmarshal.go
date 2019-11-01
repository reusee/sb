package sb

import (
	"encoding"
	gotoken "go/token"
	"math"
	"reflect"
)

func Unmarshal(stream Stream, target any) error {
	return UnmarshalValue(stream, reflect.ValueOf(target))
}

type Detokenizer interface {
	DetokenizeSB(stream Stream) error
}

func UnmarshalValue(stream Stream, ptr reflect.Value) error {

	if ptr.IsValid() {
		i := ptr.Interface()
		if v, ok := i.(Detokenizer); ok {
			return v.DetokenizeSB(stream)
		} else if v, ok := i.(encoding.BinaryUnmarshaler); ok {
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
			if err = v.UnmarshalBinary(
				[]byte(token.Value.(string)),
			); err != nil {
				return err
			}
			return nil
		} else if v, ok := i.(encoding.TextUnmarshaler); ok {
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
			if err = v.UnmarshalText(
				[]byte(token.Value.(string)),
			); err != nil {
				return err
			}
			return nil
		}
	}

	token, err := stream.Next()
	if err != nil {
		return err
	}
	if token == nil {
		return UnmarshalError{ExpectingValue}
	}

	switch token.Kind {
	case KindNil:
		return nil
	case KindArrayEnd, KindObjectEnd:
		return UnmarshalError{ExpectingValue}
	}

	valueType := ptr.Type().Elem()
	valueKind := valueType.Kind()
	hasConcreteType := false
	if valueKind != reflect.Interface {
		hasConcreteType = true
		if ptr.IsNil() {
			ptr = reflect.New(valueType)
		}
	}

	switch token.Kind {

	case KindBool:
		if hasConcreteType {
			if valueKind != reflect.Bool {
				return UnmarshalError{ExpectingBool}
			}
			ptr.Elem().SetBool(token.Value.(bool))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(bool)))
		}

	case KindInt:
		if hasConcreteType {
			if valueKind != reflect.Int {
				return UnmarshalError{ExpectingInt}
			}
			ptr.Elem().SetInt(int64(token.Value.(int)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int)))
		}

	case KindInt8:
		if hasConcreteType {
			if valueKind != reflect.Int8 {
				return UnmarshalError{ExpectingInt8}
			}
			ptr.Elem().SetInt(int64(token.Value.(int8)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int8)))
		}

	case KindInt16:
		if hasConcreteType {
			if valueKind != reflect.Int16 {
				return UnmarshalError{ExpectingInt16}
			}
			ptr.Elem().SetInt(int64(token.Value.(int16)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int16)))
		}

	case KindInt32:
		if hasConcreteType {
			if valueKind != reflect.Int32 {
				return UnmarshalError{ExpectingInt32}
			}
			ptr.Elem().SetInt(int64(token.Value.(int32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int32)))
		}

	case KindInt64:
		if hasConcreteType {
			if valueKind != reflect.Int64 {
				return UnmarshalError{ExpectingInt64}
			}
			ptr.Elem().SetInt(token.Value.(int64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(int64)))
		}

	case KindUint:
		if hasConcreteType {
			if valueKind != reflect.Uint {
				return UnmarshalError{ExpectingUint}
			}
			ptr.Elem().SetUint(uint64(token.Value.(uint)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint)))
		}

	case KindUint8:
		if hasConcreteType {
			if valueKind != reflect.Uint8 {
				return UnmarshalError{ExpectingUint8}
			}
			ptr.Elem().SetUint(uint64(token.Value.(uint8)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
		}

	case KindUint16:
		if hasConcreteType {
			if valueKind != reflect.Uint16 {
				return UnmarshalError{ExpectingUint16}
			}
			ptr.Elem().SetUint(uint64(token.Value.(uint16)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
		}

	case KindUint32:
		if hasConcreteType {
			if valueKind != reflect.Uint32 {
				return UnmarshalError{ExpectingUint32}
			}
			ptr.Elem().SetUint(uint64(token.Value.(uint32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
		}

	case KindUint64:
		if hasConcreteType {
			if valueKind != reflect.Uint64 {
				return UnmarshalError{ExpectingUint64}
			}
			ptr.Elem().SetUint(token.Value.(uint64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
		}

	case KindFloat32:
		if hasConcreteType {
			if valueKind != reflect.Float32 {
				return UnmarshalError{ExpectingFloat32}
			}
			ptr.Elem().SetFloat(float64(token.Value.(float32)))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float32)))
		}

	case KindFloat64:
		if hasConcreteType {
			if valueKind != reflect.Float64 {
				return UnmarshalError{ExpectingFloat64}
			}
			ptr.Elem().SetFloat(token.Value.(float64))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(float64)))
		}

	case KindNaN:
		if hasConcreteType {
			if valueKind != reflect.Float32 && valueKind != reflect.Float64 {
				return UnmarshalError{ExpectingFloat}
			}
			ptr.Elem().SetFloat(math.NaN())
		} else {
			ptr.Elem().Set(reflect.ValueOf(math.NaN()))
		}

	case KindString:
		if hasConcreteType {
			if valueKind != reflect.String {
				return UnmarshalError{ExpectingString}
			}
			ptr.Elem().SetString(token.Value.(string))
		} else {
			ptr.Elem().Set(reflect.ValueOf(token.Value.(string)))
		}

	case KindArray:
		if hasConcreteType {

			if valueKind == reflect.Array {
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
					if idx >= ptr.Elem().Len() {
						return UnmarshalError{TooManyElement}
					}
					err = UnmarshalValue(stream, ptr.Elem().Index(idx).Addr())
					if err != nil {
						return err
					}
					idx++
				}

			} else if valueKind == reflect.Slice {
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

			} else {
				return UnmarshalError{ExpectingSequence}
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
			if valueKind != reflect.Struct {
				return UnmarshalError{ExpectingStruct}
			}

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
			names := make(map[string]struct{})
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
				if !gotoken.IsIdentifier(name) || !gotoken.IsExported(name) {
					return UnmarshalError{BadFieldName}
				}
				if _, ok := names[name]; ok {
					return UnmarshalError{DuplicatedFieldName}
				}
				names[name] = struct{}{}

				// value
				var value any
				err = UnmarshalValue(stream, reflect.ValueOf(&value))
				if err != nil {
					return err
				}
				if value == nil {
					return UnmarshalError{ExpectingValue}
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

	case KindMap:
		if hasConcreteType {
			if valueKind != reflect.Map {
				return UnmarshalError{ExpectingMap}
			}

			keyType := valueType.Key()
			elemType := valueType.Elem()
			for {
				p, err := stream.Peek()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindMapEnd {
					stream.Next()
					break
				}
				// key
				key := reflect.New(keyType)
				if err := UnmarshalValue(stream, key); err != nil {
					return err
				}
				// value
				value := reflect.New(elemType)
				if err := UnmarshalValue(stream, value); err != nil {
					return err
				}
				if ptr.Elem().IsNil() {
					ptr.Elem().Set(reflect.MakeMap(valueType))
				}
				ptr.Elem().SetMapIndex(
					key.Elem(),
					value.Elem(),
				)
			}

		} else {
			// map[any]any
			m := make(map[any]any)
			for {
				p, err := stream.Peek()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindMapEnd {
					stream.Next()
					break
				}
				// key
				var key any
				if err := UnmarshalValue(stream, reflect.ValueOf(&key)); err != nil {
					return err
				}
				// value
				var value any
				if err := UnmarshalValue(stream, reflect.ValueOf(&value)); err != nil {
					return err
				}
				m[key] = value
			}
			ptr.Elem().Set(reflect.ValueOf(m))

		}

	}

	return nil
}
