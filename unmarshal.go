package sb

import (
	"encoding"
	gotoken "go/token"
	"math"
	"reflect"
)

type SBUnmarshaler interface {
	UnmarshalSB(vu ValueUnmarshalFunc, cont Sink) Sink
}

func Unmarshal(target any) Sink {
	return FilterSink(
		UnmarshalValue(
			UnmarshalValue,
			reflect.ValueOf(target),
			nil,
		),
		func(token *Token) bool {
			return token.Kind == KindPostTag
		},
	)
}

type ValueUnmarshalFunc func(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink

func UnmarshalValue(vu ValueUnmarshalFunc, target reflect.Value, cont Sink) Sink {
	if vu == nil {
		vu = UnmarshalValue
	}

	if target.IsValid() {
		i := target.Interface()
		if v, ok := i.(SBUnmarshaler); ok {
			return v.UnmarshalSB(vu, cont)

		} else if v, ok := i.(encoding.BinaryUnmarshaler); ok {
			return func(p *Token) (Sink, error) {
				if p == nil {
					return nil, UnmarshalError{ExpectingString}
				}
				token := *p
				if token.Kind != KindString {
					return nil, UnmarshalError{ExpectingString}
				}
				if err := v.UnmarshalBinary(
					[]byte(token.Value.(string)),
				); err != nil {
					return nil, err
				}
				return cont, nil
			}

		} else if v, ok := i.(encoding.TextUnmarshaler); ok {
			return func(p *Token) (Sink, error) {
				if p == nil {
					return nil, UnmarshalError{ExpectingString}
				}
				token := *p
				if token.Kind != KindString {
					return nil, UnmarshalError{ExpectingString}
				}
				if err := v.UnmarshalText(
					[]byte(token.Value.(string)),
				); err != nil {
					return nil, err
				}
				return cont, nil
			}
		}

	}

	return func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingValue}
		}

		switch token.Kind {
		case KindNil:
			return cont, nil
		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			return nil, UnmarshalError{ExpectingValue}
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
			return nil, UnmarshalError{BadTargetType}
		}

		hasConcreteType := false
		if valueKind == reflect.Ptr {
			// deref
			t := reflect.New(valueType.Elem())
			return vu(
				vu,
				t,
				func(token *Token) (Sink, error) {
					// target will not set unless no error
					target.Elem().Set(t)
					if cont != nil {
						return cont(token)
					}
					return nil, nil
				},
			)(token)
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
					return nil, UnmarshalError{ExpectingBool}
				}
				target.Elem().SetBool(token.Value.(bool))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(bool)))
			}

		case KindInt:
			if hasConcreteType {
				if valueKind != reflect.Int {
					return nil, UnmarshalError{ExpectingInt}
				}
				target.Elem().SetInt(int64(token.Value.(int)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int)))
			}

		case KindInt8:
			if hasConcreteType {
				if valueKind != reflect.Int8 {
					return nil, UnmarshalError{ExpectingInt8}
				}
				target.Elem().SetInt(int64(token.Value.(int8)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int8)))
			}

		case KindInt16:
			if hasConcreteType {
				if valueKind != reflect.Int16 {
					return nil, UnmarshalError{ExpectingInt16}
				}
				target.Elem().SetInt(int64(token.Value.(int16)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int16)))
			}

		case KindInt32:
			if hasConcreteType {
				if valueKind != reflect.Int32 {
					return nil, UnmarshalError{ExpectingInt32}
				}
				target.Elem().SetInt(int64(token.Value.(int32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int32)))
			}

		case KindInt64:
			if hasConcreteType {
				if valueKind != reflect.Int64 {
					return nil, UnmarshalError{ExpectingInt64}
				}
				target.Elem().SetInt(token.Value.(int64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int64)))
			}

		case KindUint:
			if hasConcreteType {
				if valueKind != reflect.Uint {
					return nil, UnmarshalError{ExpectingUint}
				}
				target.Elem().SetUint(uint64(token.Value.(uint)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint)))
			}

		case KindUint8:
			if hasConcreteType {
				if valueKind != reflect.Uint8 {
					return nil, UnmarshalError{ExpectingUint8}
				}
				target.Elem().SetUint(uint64(token.Value.(uint8)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
			}

		case KindUint16:
			if hasConcreteType {
				if valueKind != reflect.Uint16 {
					return nil, UnmarshalError{ExpectingUint16}
				}
				target.Elem().SetUint(uint64(token.Value.(uint16)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
			}

		case KindUint32:
			if hasConcreteType {
				if valueKind != reflect.Uint32 {
					return nil, UnmarshalError{ExpectingUint32}
				}
				target.Elem().SetUint(uint64(token.Value.(uint32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
			}

		case KindUint64:
			if hasConcreteType {
				if valueKind != reflect.Uint64 {
					return nil, UnmarshalError{ExpectingUint64}
				}
				target.Elem().SetUint(token.Value.(uint64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
			}

		case KindFloat32:
			if hasConcreteType {
				if valueKind != reflect.Float32 {
					return nil, UnmarshalError{ExpectingFloat32}
				}
				target.Elem().SetFloat(float64(token.Value.(float32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(float32)))
			}

		case KindFloat64:
			if hasConcreteType {
				if valueKind != reflect.Float64 {
					return nil, UnmarshalError{ExpectingFloat64}
				}
				target.Elem().SetFloat(token.Value.(float64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(float64)))
			}

		case KindNaN:
			if hasConcreteType {
				if valueKind != reflect.Float32 && valueKind != reflect.Float64 {
					return nil, UnmarshalError{ExpectingFloat}
				}
				target.Elem().SetFloat(math.NaN())
			} else {
				target.Elem().Set(reflect.ValueOf(math.NaN()))
			}

		case KindString:
			if hasConcreteType {
				if valueKind != reflect.String {
					return nil, UnmarshalError{ExpectingString}
				}
				target.Elem().SetString(token.Value.(string))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(string)))
			}

		case KindBytes:
			if hasConcreteType {
				if !isBytes(valueType) {
					return nil, UnmarshalError{ExpectingBytes}
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
					return UnmarshalArray(
						vu,
						target,
						cont,
					)(token)

				} else if valueKind == reflect.Slice {
					// slice
					return UnmarshalSlice(
						vu,
						target,
						valueType,
						cont,
					)(token)

				} else {
					return nil, UnmarshalError{ExpectingSequence}
				}

			} else {
				// generic slice
				return UnmarshalGenericSlice(
					vu,
					target,
					cont,
				)(token)
			}

		case KindObject:
			if hasConcreteType {
				if valueKind != reflect.Struct {
					return nil, UnmarshalError{ExpectingStruct}
				}
				return UnmarshalStruct(
					vu,
					target,
					valueType,
					cont,
				)(token)

			} else {
				// construct new type
				return UnmarshalNewStruct(
					vu,
					target,
					cont,
				)(token)
			}

		case KindMap:
			if hasConcreteType {
				if valueKind != reflect.Map {
					return nil, UnmarshalError{ExpectingMap}
				}
				return UnmarshalMap(
					vu,
					target,
					valueType,
					cont,
				)(token)

			} else {
				// map[any]any
				return UnmarshalGenericMap(
					vu,
					target,
					cont,
				)(token)

			}

		case KindTuple:
			if hasConcreteType {
				if valueKind != reflect.Func {
					return nil, UnmarshalError{ExpectingTuple}
				}

				numOut := valueType.NumOut()
				numIn := valueType.NumIn()

				if numOut > 0 && numIn == 0 {
					// func() (...)
					return UnmarshalTuple(
						vu,
						target,
						valueType,
						cont,
					)(token)

				} else if numIn > 0 &&
					(numOut == 0 ||
						(numOut == 1 && valueType.Out(0) == errorType)) {
					// func(...) or func(...) error

					if numIn == 1 && valueType.IsVariadic() {
						return UnmarshalTupleCall(
							vu,
							target,
							valueType,
							cont,
						)(token)

					} else {
						return UnmarshalTupleCallErr(
							vu,
							target,
							valueType,
							cont,
						)(token)
					}

				} else {
					return nil, UnmarshalError{BadTupleType}
				}

			} else {
				// func() (...any)
				return UnmarshalTupleVariadic(
					vu,
					target,
					cont,
				)(token)
			}

		default:
			return nil, UnmarshalError{BadTokenKind}
		}

		return cont, nil
	}

}

func UnmarshalArray(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink {
	return ExpectKind(
		KindArray,
		unmarshalArray(vu, target, 0, cont),
	)
}

func unmarshalArray(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	idx int,
	cont Sink,
) Sink {

	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindArrayEnd {
			return cont, nil
		}
		if idx >= target.Elem().Len() {
			return nil, UnmarshalError{TooManyElement}
		}

		e := target.Elem().Index(idx).Addr()
		idx++
		return vu(
			vu,
			e,
			sink,
		)(p)

	}
	return sink
}

func UnmarshalSlice(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	slice := target.Elem()
	return ExpectKind(
		KindArray,
		unmarshalSlice(
			vu,
			target, valueType,
			slice, cont,
		),
	)
}

func unmarshalSlice(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	slice reflect.Value,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindArrayEnd {
			target.Elem().Set(slice)
			return cont, nil
		}
		elemPtr := reflect.New(valueType.Elem())
		slice = reflect.Append(slice, elemPtr.Elem())

		return vu(
			vu,
			slice.Index(slice.Len()-1).Addr(),
			sink,
		)(p)

	}
	return sink
}

func UnmarshalGenericSlice(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink {
	var slice []any
	return ExpectKind(
		KindArray,
		unmarshalGenericSlice(
			vu,
			target, slice, cont,
		),
	)
}

func unmarshalGenericSlice(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	slice []any,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindArrayEnd {
			target.Elem().Set(reflect.ValueOf(slice))
			return cont, nil
		}

		var value any
		return vu(
			vu,
			reflect.ValueOf(&value),
			func(token *Token) (Sink, error) {
				slice = append(slice, value)
				return sink(token)
			},
		)(p)

	}
	return sink
}

func UnmarshalStruct(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	return ExpectKind(
		KindObject,
		unmarshalStruct(
			vu,
			target, valueType, cont,
		),
	)
}

func unmarshalStruct(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindObjectEnd {
			return cont, nil
		}
		var name string

		return vu(
			vu,
			reflect.ValueOf(&name),
			func(token *Token) (Sink, error) {
				field, ok := valueType.FieldByName(name)
				if !ok {
					// skip next value
					var value any
					return vu(
						vu,
						reflect.ValueOf(&value),
						sink,
					)(token)

				} else {
					return vu(
						vu,
						target.Elem().FieldByIndex(field.Index).Addr(),
						sink,
					)(token)
				}

			},
		)(p)

	}
	return sink
}

func UnmarshalNewStruct(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink {
	var values []any
	var fields []reflect.StructField
	names := make(map[string]struct{})
	return ExpectKind(
		KindObject,
		unmarshalNewStruct(
			vu,
			target, values, fields, names, cont,
		),
	)
}

func unmarshalNewStruct(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	values []any,
	fields []reflect.StructField,
	names map[string]struct{},
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindObjectEnd {
			structType := reflect.StructOf(fields)
			structPtr := reflect.New(structType)
			for i, value := range values {
				structPtr.Elem().Field(i).Set(reflect.ValueOf(value))
			}
			target.Elem().Set(structPtr.Elem())
			return cont, nil
		}

		var name string
		return vu(
			vu,
			reflect.ValueOf(&name),
			func(token *Token) (Sink, error) {
				if !gotoken.IsIdentifier(name) || !gotoken.IsExported(name) {
					return nil, UnmarshalError{BadFieldName}
				}
				if _, ok := names[name]; ok {
					return nil, UnmarshalError{DuplicatedFieldName}
				}
				names[name] = struct{}{}
				var value any

				return vu(
					vu,
					reflect.ValueOf(&value),
					func(token *Token) (Sink, error) {
						if value == nil {
							return nil, UnmarshalError{ExpectingValue}
						}
						values = append(values, value)
						fields = append(fields, reflect.StructField{
							Name: name,
							Type: reflect.TypeOf(value),
						})

						return sink(token)

					},
				)(token)

			},
		)(p)

	}
	return sink
}

func UnmarshalMap(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	keyType := valueType.Key()
	elemType := valueType.Elem()
	return ExpectKind(
		KindMap,
		unmarshalMap(
			vu,
			target, valueType, keyType, elemType, cont,
		),
	)
}

func unmarshalMap(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	keyType reflect.Type,
	elemType reflect.Type,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindMapEnd {
			return cont, nil
		}

		key := reflect.New(keyType)
		return vu(
			vu,
			key,
			func(token *Token) (Sink, error) {
				value := reflect.New(elemType)

				return vu(
					vu,
					value,
					func(token *Token) (Sink, error) {
						if target.Elem().IsNil() {
							target.Elem().Set(reflect.MakeMap(valueType))
						}
						target.Elem().SetMapIndex(
							key.Elem(),
							value.Elem(),
						)
						return sink(token)
					},
				)(token)

			},
		)(p)
	}
	return sink
}

func UnmarshalGenericMap(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink {
	m := make(map[any]any)
	return ExpectKind(
		KindMap,
		unmarshalGenericMap(
			vu,
			target, m, cont,
		),
	)
}

func unmarshalGenericMap(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	m map[any]any,
	cont Sink,
) Sink {
	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindMapEnd {
			target.Elem().Set(reflect.ValueOf(m))
			return cont, nil
		}

		var key any
		return vu(
			vu,
			reflect.ValueOf(&key),
			func(token *Token) (Sink, error) {
				if key == nil {
					return nil, UnmarshalError{ExpectingValue}
				} else if !reflect.TypeOf(key).Comparable() {
					return nil, UnmarshalError{BadMapKey}
				} else if f, ok := key.(float64); ok && math.IsNaN(f) {
					return nil, UnmarshalError{BadMapKey}
				} else if f, ok := key.(float32); ok && math.IsNaN(float64(f)) {
					return nil, UnmarshalError{BadMapKey}
				}
				var value any

				return vu(
					vu,
					reflect.ValueOf(&value),
					func(token *Token) (Sink, error) {
						m[key] = value
						return sink(token)
					},
				)(token)

			},
		)(p)
	}
	return sink
}

func UnmarshalTuple(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	numOut := valueType.NumOut()
	var items []reflect.Value
	return ExpectKind(
		KindTuple,
		unmarshalTuple(
			vu,
			target, valueType, numOut, 0, items, cont,
		),
	)
}

func unmarshalTuple(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	numOut int,
	i int,
	items []reflect.Value,
	cont Sink,
) Sink {

	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if i >= numOut {
			if p.Kind != KindTupleEnd {
				return nil, UnmarshalError{TooManyElement}
			}
			target.Elem().Set(reflect.MakeFunc(
				valueType,
				func(args []reflect.Value) []reflect.Value {
					return items
				},
			))
			return cont, nil
		}

		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindTupleEnd {
			return nil, UnmarshalError{ExpectingValue}
		}
		itemType := valueType.Out(i)
		value := reflect.New(itemType)
		return vu(
			vu,
			value,
			func(token *Token) (Sink, error) {
				items = append(items, value.Elem())
				i++
				return sink(token)
			},
		)(p)
	}
	return sink
}

func UnmarshalTupleCall(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	itemType := valueType.In(0).Elem()
	var items []reflect.Value
	return ExpectKind(
		KindTuple,
		unmarshalTupleCall(
			vu,
			target, itemType, items, cont,
		),
	)
}

func unmarshalTupleCall(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	itemType reflect.Type,
	items []reflect.Value,
	cont Sink,
) Sink {

	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindTupleEnd {
			if !target.IsNil() {
				rets := target.Call(items)
				if len(rets) > 0 {
					i := rets[0].Interface()
					if i != nil {
						return nil, i.(error)
					}
					return cont, nil
				}
			}
			return cont, nil
		}

		value := reflect.New(itemType)
		return vu(
			vu,
			value,
			func(token *Token) (Sink, error) {
				items = append(items, value.Elem())
				return sink(token)
			},
		)(p)

	}
	return sink
}

func UnmarshalTupleCallErr(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	numIn := valueType.NumIn()
	var items []reflect.Value
	return ExpectKind(
		KindTuple,
		unmarshalTupleCallErr(
			vu,
			target, valueType, numIn, 0, items, cont,
		),
	)
}

func unmarshalTupleCallErr(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	valueType reflect.Type,
	numIn int,
	i int,
	items []reflect.Value,
	cont Sink,
) Sink {

	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if i >= numIn {
			if p.Kind != KindTupleEnd {
				return nil, UnmarshalError{TooManyElement}
			}
			if !target.IsNil() {
				rets := target.Call(items)
				if len(rets) > 0 {
					i := rets[0].Interface()
					if i != nil {
						return nil, i.(error)
					}
					return cont, nil
				}
			}
			return cont, nil
		}

		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindTupleEnd {
			return nil, UnmarshalError{ExpectingValue}
		}
		itemType := valueType.In(i)
		value := reflect.New(itemType)
		return vu(
			vu,
			value,
			func(token *Token) (Sink, error) {
				items = append(items, value.Elem())
				i++
				return sink(token)
			},
		)(p)
	}

	return sink
}

func UnmarshalTupleVariadic(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	cont Sink,
) Sink {
	var items []reflect.Value
	var itemTypes []reflect.Type
	return ExpectKind(
		KindTuple,
		unmarshalTupleVariadic(
			vu,
			target, items, itemTypes, cont,
		),
	)
}

func unmarshalTupleVariadic(
	vu ValueUnmarshalFunc,
	target reflect.Value,
	items []reflect.Value,
	itemTypes []reflect.Type,
	cont Sink,
) Sink {

	var sink Sink
	sink = func(p *Token) (Sink, error) {
		if p == nil {
			return nil, UnmarshalError{ExpectingValue}
		}
		if p.Kind == KindTupleEnd {
			if len(itemTypes) > 50 {
				return nil, UnmarshalError{TooManyElement}
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
			return cont, nil
		}

		var obj any
		value := reflect.ValueOf(&obj)
		itemTypes = append(itemTypes, anyType)
		return vu(
			vu,
			value,
			func(token *Token) (Sink, error) {
				items = append(items, value.Elem())
				return sink(token)
			},
		)(p)
	}

	return sink
}
