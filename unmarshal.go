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

type SBUnmarshaler interface {
	UnmarshalSB(stream Stream) error
}

func UnmarshalValue(stream Stream, target reflect.Value) error {

	stream = Filter(stream, func(token *Token) bool {
		return token.Kind == KindPostTag
	})

	if target.IsValid() {
		i := target.Interface()
		if v, ok := i.(SBUnmarshaler); ok {
			return v.UnmarshalSB(stream)
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
	case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
		return UnmarshalError{ExpectingValue}
	}

	targetType := target.Type()
	targetKind := targetType.Kind()
	var valueType reflect.Type
	var valueKind reflect.Kind
	if targetKind == reflect.Func {
		valueType = targetType
		valueKind = targetKind
	} else if targetKind == reflect.Ptr {
		valueType = targetType.Elem()
		valueKind = valueType.Kind()
	} else {
		return UnmarshalError{BadTargetType}
	}

	hasConcreteType := false
	if valueKind == reflect.Ptr {
		// deref
		if target.Elem().IsNil() {
			target.Elem().Set(reflect.New(valueType.Elem()))
		}
		return UnmarshalValue(
			&unreadToken{token, stream},
			target.Elem(),
		)
	} else if valueKind != reflect.Interface {
		hasConcreteType = true
		if targetKind == reflect.Ptr && target.IsNil() {
			target = reflect.New(valueType)
		}
	}

	switch token.Kind {

	case KindBool:
		if hasConcreteType {
			if valueKind != reflect.Bool {
				return UnmarshalError{ExpectingBool}
			}
			target.Elem().SetBool(token.Value.(bool))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(bool)))
		}

	case KindInt:
		if hasConcreteType {
			if valueKind != reflect.Int {
				return UnmarshalError{ExpectingInt}
			}
			target.Elem().SetInt(int64(token.Value.(int)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(int)))
		}

	case KindInt8:
		if hasConcreteType {
			if valueKind != reflect.Int8 {
				return UnmarshalError{ExpectingInt8}
			}
			target.Elem().SetInt(int64(token.Value.(int8)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(int8)))
		}

	case KindInt16:
		if hasConcreteType {
			if valueKind != reflect.Int16 {
				return UnmarshalError{ExpectingInt16}
			}
			target.Elem().SetInt(int64(token.Value.(int16)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(int16)))
		}

	case KindInt32:
		if hasConcreteType {
			if valueKind != reflect.Int32 {
				return UnmarshalError{ExpectingInt32}
			}
			target.Elem().SetInt(int64(token.Value.(int32)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(int32)))
		}

	case KindInt64:
		if hasConcreteType {
			if valueKind != reflect.Int64 {
				return UnmarshalError{ExpectingInt64}
			}
			target.Elem().SetInt(token.Value.(int64))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(int64)))
		}

	case KindUint:
		if hasConcreteType {
			if valueKind != reflect.Uint {
				return UnmarshalError{ExpectingUint}
			}
			target.Elem().SetUint(uint64(token.Value.(uint)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(uint)))
		}

	case KindUint8:
		if hasConcreteType {
			if valueKind != reflect.Uint8 {
				return UnmarshalError{ExpectingUint8}
			}
			target.Elem().SetUint(uint64(token.Value.(uint8)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
		}

	case KindUint16:
		if hasConcreteType {
			if valueKind != reflect.Uint16 {
				return UnmarshalError{ExpectingUint16}
			}
			target.Elem().SetUint(uint64(token.Value.(uint16)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
		}

	case KindUint32:
		if hasConcreteType {
			if valueKind != reflect.Uint32 {
				return UnmarshalError{ExpectingUint32}
			}
			target.Elem().SetUint(uint64(token.Value.(uint32)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
		}

	case KindUint64:
		if hasConcreteType {
			if valueKind != reflect.Uint64 {
				return UnmarshalError{ExpectingUint64}
			}
			target.Elem().SetUint(token.Value.(uint64))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
		}

	case KindFloat32:
		if hasConcreteType {
			if valueKind != reflect.Float32 {
				return UnmarshalError{ExpectingFloat32}
			}
			target.Elem().SetFloat(float64(token.Value.(float32)))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(float32)))
		}

	case KindFloat64:
		if hasConcreteType {
			if valueKind != reflect.Float64 {
				return UnmarshalError{ExpectingFloat64}
			}
			target.Elem().SetFloat(token.Value.(float64))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(float64)))
		}

	case KindNaN:
		if hasConcreteType {
			if valueKind != reflect.Float32 && valueKind != reflect.Float64 {
				return UnmarshalError{ExpectingFloat}
			}
			target.Elem().SetFloat(math.NaN())
		} else {
			target.Elem().Set(reflect.ValueOf(math.NaN()))
		}

	case KindString:
		if hasConcreteType {
			if valueKind != reflect.String {
				return UnmarshalError{ExpectingString}
			}
			target.Elem().SetString(token.Value.(string))
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.(string)))
		}

	case KindBytes:
		if hasConcreteType {
			if !isBytes(valueType) {
				return UnmarshalError{ExpectingBytes}
			}
			if valueKind == reflect.Slice {
				// slice
				target.Elem().Set(reflect.ValueOf(token.Value.([]byte)))
			} else {
				// array
				reflect.Copy(
					target.Elem().Slice(0, target.Elem().Len()),
					reflect.ValueOf(token.Value.([]byte)),
				)
			}
		} else {
			target.Elem().Set(reflect.ValueOf(token.Value.([]byte)))
		}

	case KindArray:
		if hasConcreteType {

			if valueKind == reflect.Array {
				// array
				idx := 0
				for {
					p, err := stream.Next()
					if err != nil {
						return err
					}
					if p == nil {
						return UnmarshalError{ExpectingValue}
					}
					if p.Kind == KindArrayEnd {
						break
					}
					if idx >= target.Elem().Len() {
						return UnmarshalError{TooManyElement}
					}
					err = UnmarshalValue(
						&unreadToken{p, stream},
						target.Elem().Index(idx).Addr(),
					)
					if err != nil {
						return err
					}
					idx++
				}

			} else if valueKind == reflect.Slice {
				// slice
				slice := target.Elem()
				for {
					p, err := stream.Next()
					if err != nil {
						return err
					}
					if p == nil {
						return UnmarshalError{ExpectingValue}
					}
					if p.Kind == KindArrayEnd {
						break
					}
					elemPtr := reflect.New(valueType.Elem())
					err = UnmarshalValue(
						&unreadToken{p, stream},
						elemPtr,
					)
					if err != nil {
						return err
					}
					slice = reflect.Append(slice, elemPtr.Elem())
				}
				target.Elem().Set(slice)

			} else {
				return UnmarshalError{ExpectingSequence}
			}

		} else {
			// generic slice
			var slice []any
			for {
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindArrayEnd {
					break
				}
				var elem any
				err = UnmarshalValue(
					&unreadToken{p, stream},
					reflect.ValueOf(&elem),
				)
				if err != nil {
					return err
				}
				slice = append(slice, elem)
			}
			target.Elem().Set(reflect.ValueOf(slice))
		}

	case KindObject:
		if hasConcreteType {
			if valueKind != reflect.Struct {
				return UnmarshalError{ExpectingStruct}
			}

			for {
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindObjectEnd {
					break
				}

				// name
				var name string
				err = UnmarshalValue(
					&unreadToken{p, stream},
					reflect.ValueOf(&name),
				)
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
				err = UnmarshalValue(stream, target.Elem().FieldByIndex(field.Index).Addr())
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
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindObjectEnd {
					break
				}

				// name
				var name string
				err = UnmarshalValue(
					&unreadToken{p, stream},
					reflect.ValueOf(&name),
				)
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
			target.Elem().Set(structPtr.Elem())
		}

	case KindMap:
		if hasConcreteType {
			if valueKind != reflect.Map {
				return UnmarshalError{ExpectingMap}
			}

			keyType := valueType.Key()
			elemType := valueType.Elem()
			for {
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindMapEnd {
					break
				}
				// key
				key := reflect.New(keyType)
				if err := UnmarshalValue(
					&unreadToken{p, stream},
					key,
				); err != nil {
					return err
				}
				// value
				value := reflect.New(elemType)
				if err := UnmarshalValue(stream, value); err != nil {
					return err
				}
				if target.Elem().IsNil() {
					target.Elem().Set(reflect.MakeMap(valueType))
				}
				target.Elem().SetMapIndex(
					key.Elem(),
					value.Elem(),
				)
			}

		} else {
			// map[any]any
			m := make(map[any]any)
			for {
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindMapEnd {
					break
				}
				// key
				var key any
				if err := UnmarshalValue(
					&unreadToken{p, stream},
					reflect.ValueOf(&key),
				); err != nil {
					return err
				}
				if key == nil {
					return UnmarshalError{ExpectingValue}
				} else if !reflect.TypeOf(key).Comparable() {
					return UnmarshalError{BadMapKey}
				} else if f, ok := key.(float64); ok && math.IsNaN(f) {
					return UnmarshalError{BadMapKey}
				} else if f, ok := key.(float32); ok && math.IsNaN(float64(f)) {
					return UnmarshalError{BadMapKey}
				}
				// value
				var value any
				if err := UnmarshalValue(stream, reflect.ValueOf(&value)); err != nil {
					return err
				}
				m[key] = value
			}
			target.Elem().Set(reflect.ValueOf(m))

		}

	case KindTuple:
		if hasConcreteType {
			if valueKind != reflect.Func {
				return UnmarshalError{ExpectingTuple}
			}

			numOut := valueType.NumOut()
			numIn := valueType.NumIn()

			if numOut > 0 && numIn == 0 {
				// func() (...)
				var items []reflect.Value
				for i := 0; i < numOut; i++ {
					p, err := stream.Next()
					if err != nil {
						return err
					}
					if p == nil {
						return UnmarshalError{ExpectingValue}
					}
					if p.Kind == KindTupleEnd {
						return UnmarshalError{ExpectingValue}
					}
					itemType := valueType.Out(i)
					value := reflect.New(itemType)
					if err := UnmarshalValue(
						&unreadToken{p, stream},
						value,
					); err != nil {
						return err
					}
					items = append(items, value.Elem())
				}
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p.Kind != KindTupleEnd {
					return UnmarshalError{TooManyElement}
				}
				target.Elem().Set(reflect.MakeFunc(
					valueType,
					func(args []reflect.Value) []reflect.Value {
						return items
					},
				))

			} else if numIn > 0 &&
				(numOut == 0 ||
					(numOut == 1 && valueType.Out(0) == errorType)) {
				// func(...) or func(...) error

				if numIn == 1 && valueType.IsVariadic() {
					itemType := valueType.In(0).Elem()
					var items []reflect.Value
					for {
						p, err := stream.Next()
						if err != nil {
							return err
						}
						if p == nil {
							return UnmarshalError{ExpectingValue}
						}
						if p.Kind == KindTupleEnd {
							break
						}
						value := reflect.New(itemType)
						if err := UnmarshalValue(
							&unreadToken{p, stream},
							value,
						); err != nil {
							return err
						}
						items = append(items, value.Elem())
					}
					if !target.IsNil() {
						rets := target.Call(items)
						if len(rets) > 0 {
							i := rets[0].Interface()
							if i != nil {
								return i.(error)
							}
							return nil
						}
					}

				} else {
					var items []reflect.Value
					for i := 0; i < numIn; i++ {
						p, err := stream.Next()
						if err != nil {
							return err
						}
						if p == nil {
							return UnmarshalError{ExpectingValue}
						}
						if p.Kind == KindTupleEnd {
							return UnmarshalError{ExpectingValue}
						}
						itemType := valueType.In(i)
						value := reflect.New(itemType)
						if err := UnmarshalValue(
							&unreadToken{p, stream},
							value,
						); err != nil {
							return err
						}
						items = append(items, value.Elem())
					}
					p, err := stream.Next()
					if err != nil {
						return err
					}
					if p.Kind != KindTupleEnd {
						return UnmarshalError{TooManyElement}
					}
					if !target.IsNil() {
						rets := target.Call(items)
						if len(rets) > 0 {
							i := rets[0].Interface()
							if i != nil {
								return i.(error)
							}
							return nil
						}
					}
				}

			} else {
				return UnmarshalError{BadTupleType}
			}

		} else {
			// func() (...any)
			var items []reflect.Value
			var itemTypes []reflect.Type
			for {
				p, err := stream.Next()
				if err != nil {
					return err
				}
				if p == nil {
					return UnmarshalError{ExpectingValue}
				}
				if p.Kind == KindTupleEnd {
					break
				}
				var obj any
				value := reflect.ValueOf(&obj)
				itemTypes = append(itemTypes, anyType)
				if err := UnmarshalValue(
					&unreadToken{p, stream},
					value,
				); err != nil {
					return err
				}
				items = append(items, value.Elem())
			}
			if len(itemTypes) > 50 {
				return UnmarshalError{TooManyElement}
			}
			target.Elem().Set(reflect.MakeFunc(
				reflect.FuncOf(
					[]reflect.Type{},
					itemTypes,
					false,
				),
				func(args []reflect.Value) []reflect.Value {
					return items
				},
			))
		}

	default:
		return UnmarshalError{BadTokenKind}
	}

	return nil
}

type unreadToken struct {
	token  *Token
	stream Stream
}

var _ Stream = new(unreadToken)

func (s *unreadToken) Next() (token *Token, err error) {
	if s.token != nil {
		token = s.token
		s.token = nil
		return
	}
	return s.stream.Next()
}
