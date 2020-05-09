package sb

import (
	"encoding"
	"fmt"
	gotoken "go/token"
	"math"
	"reflect"
)

type SBUnmarshaler interface {
	UnmarshalSB(ctx Ctx, cont Sink) Sink
}

func Unmarshal(target any) Sink {
	return UnmarshalValue(
		Ctx{
			Unmarshal: UnmarshalValue,
		},
		reflect.ValueOf(target),
		nil,
	)
}

func UnmarshalValue(ctx Ctx, target reflect.Value, cont Sink) Sink {
	if ctx.Unmarshal == nil {
		ctx.Unmarshal = UnmarshalValue
	}

	return func(token *Token) (Sink, error) {

		if target.IsValid() {
			switch v := target.Interface().(type) {

			case *int:
				if token.Kind == KindInt && v != nil {
					*v = token.Value.(int)
					return cont, nil
				}

			case *string:
				if token.Kind == KindString && v != nil {
					*v = token.Value.(string)
					return cont, nil
				}

			case *bool:
				if token.Kind == KindBool && v != nil {
					*v = token.Value.(bool)
					return cont, nil
				}

			case *uint32:
				if token.Kind == KindUint32 && v != nil {
					*v = token.Value.(uint32)
					return cont, nil
				}

			case *int64:
				if token.Kind == KindInt64 && v != nil {
					*v = token.Value.(int64)
					return cont, nil
				}

			case *uint64:
				if token.Kind == KindUint64 && v != nil {
					*v = token.Value.(uint64)
					return cont, nil
				}

			case *uint16:
				if token.Kind == KindUint16 && v != nil {
					*v = token.Value.(uint16)
					return cont, nil
				}

			case *uint8:
				if token.Kind == KindUint8 && v != nil {
					*v = token.Value.(uint8)
					return cont, nil
				}

			case *int32:
				if token.Kind == KindInt32 && v != nil {
					*v = token.Value.(int32)
					return cont, nil
				}

			case *uint:
				if token.Kind == KindUint && v != nil {
					*v = token.Value.(uint)
					return cont, nil
				}

			case *float64:
				if token.Kind == KindFloat64 && v != nil {
					*v = token.Value.(float64)
					return cont, nil
				}

			case *int8:
				if token.Kind == KindInt8 && v != nil {
					*v = token.Value.(int8)
					return cont, nil
				}

			case *float32:
				if token.Kind == KindFloat32 && v != nil {
					*v = token.Value.(float32)
					return cont, nil
				}

			case *int16:
				if token.Kind == KindInt16 && v != nil {
					*v = token.Value.(int16)
					return cont, nil
				}

			case SBUnmarshaler:
				return v.UnmarshalSB(ctx, cont)(token)

			case encoding.BinaryUnmarshaler:
				if token == nil {
					return nil, UnmarshalError{ExpectingString}
				}
				if token.Kind != KindString {
					return nil, UnmarshalError{ExpectingString}
				}
				if err := v.UnmarshalBinary(
					[]byte(token.Value.(string)),
				); err != nil {
					return nil, err
				}
				return cont, nil

			case encoding.TextUnmarshaler:
				if token == nil {
					return nil, UnmarshalError{ExpectingString}
				}
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
			return ctx.Unmarshal(
				ctx,
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
						ctx,
						target,
						cont,
					)(token)

				} else if valueKind == reflect.Slice {
					// slice
					return UnmarshalSlice(
						ctx,
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
					ctx,
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
					ctx,
					target,
					valueType,
					cont,
				)(token)

			} else {
				// construct new type
				return UnmarshalNewStruct(
					ctx,
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
					ctx,
					target,
					valueType,
					cont,
				)(token)

			} else {
				// map[any]any
				return UnmarshalGenericMap(
					ctx,
					target,
					cont,
				)(token)

			}

		case KindTuple:
			if hasConcreteType {
				if valueKind != reflect.Func {
					return nil, UnmarshalError{ExpectingTuple}
				}
			}
			return UnmarshalTuple(
				ctx,
				target,
				valueType,
				cont,
			)(token)

		default:
			return nil, UnmarshalError{BadTokenKind}
		}

		return cont, nil
	}

}

func UnmarshalArray(
	ctx Ctx,
	target reflect.Value,
	cont Sink,
) Sink {
	return ExpectKind(
		KindArray,
		unmarshalArray(ctx, target, 0, cont),
	)
}

func unmarshalArray(
	ctx Ctx,
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
		return ctx.Unmarshal(
			ctx,
			e,
			sink,
		)(p)

	}
	return sink
}

func UnmarshalSlice(
	ctx Ctx,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	slice := target.Elem()
	return ExpectKind(
		KindArray,
		unmarshalSlice(
			ctx,
			target, valueType,
			slice, cont,
		),
	)
}

func unmarshalSlice(
	ctx Ctx,
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

		return ctx.Unmarshal(
			ctx,
			slice.Index(slice.Len()-1).Addr(),
			sink,
		)(p)

	}
	return sink
}

func UnmarshalGenericSlice(
	ctx Ctx,
	target reflect.Value,
	cont Sink,
) Sink {
	var slice []any
	return ExpectKind(
		KindArray,
		unmarshalGenericSlice(
			ctx,
			target, slice, cont,
		),
	)
}

func unmarshalGenericSlice(
	ctx Ctx,
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
		return ctx.Unmarshal(
			ctx,
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
	ctx Ctx,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	return ExpectKind(
		KindObject,
		unmarshalStruct(
			ctx,
			target, valueType, cont,
		),
	)
}

func unmarshalStruct(
	ctx Ctx,
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

		return ctx.Unmarshal(
			ctx,
			reflect.ValueOf(&name),
			func(token *Token) (Sink, error) {
				field, ok := valueType.FieldByName(name)
				if !ok {
					if ctx.DisallowUnknownStructFields {
						return nil, fmt.Errorf("field %s: %w", name, UnmarshalError{UnknownFieldName})
					} else {
						// skip next value
						var value any
						return ctx.Unmarshal(
							ctx,
							reflect.ValueOf(&value),
							sink,
						)(token)
					}

				} else {
					return ctx.Unmarshal(
						ctx,
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
	ctx Ctx,
	target reflect.Value,
	cont Sink,
) Sink {
	var values []any
	var fields []reflect.StructField
	names := make(map[string]struct{})
	return ExpectKind(
		KindObject,
		unmarshalNewStruct(
			ctx,
			target, values, fields, names, cont,
		),
	)
}

func unmarshalNewStruct(
	ctx Ctx,
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
		return ctx.Unmarshal(
			ctx,
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

				return ctx.Unmarshal(
					ctx,
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
	ctx Ctx,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {
	keyType := valueType.Key()
	elemType := valueType.Elem()
	return ExpectKind(
		KindMap,
		unmarshalMap(
			ctx,
			target, valueType, keyType, elemType, cont,
		),
	)
}

func unmarshalMap(
	ctx Ctx,
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
		return ctx.Unmarshal(
			ctx,
			key,
			func(token *Token) (Sink, error) {
				value := reflect.New(elemType)

				return ctx.Unmarshal(
					ctx,
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
	ctx Ctx,
	target reflect.Value,
	cont Sink,
) Sink {
	m := make(map[any]any)
	return ExpectKind(
		KindMap,
		unmarshalGenericMap(
			ctx,
			target, m, cont,
		),
	)
}

func unmarshalGenericMap(
	ctx Ctx,
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
		return ctx.Unmarshal(
			ctx,
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

				return ctx.Unmarshal(
					ctx,
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

var ellipsesType = reflect.TypeOf((*[]any)(nil)).Elem()

func UnmarshalTuple(
	ctx Ctx,
	target reflect.Value,
	valueType reflect.Type,
	cont Sink,
) Sink {

	var concreteTypes []reflect.Type
	if valueType.Kind() == reflect.Func {
		numIn := valueType.NumIn()
		numOut := valueType.NumOut()
		if numIn == 0 {
			// return only
			for i := 0; i < numOut; i++ {
				t := valueType.Out(i)
				concreteTypes = append(concreteTypes, t)
			}
		} else {
			for i := 0; i < numIn; i++ {
				t := valueType.In(i)
				if t == ellipsesType {
					continue
				}
				concreteTypes = append(concreteTypes, t)
			}
		}
	}

	return ExpectKind(
		KindTuple,
		unmarshalTuple(
			ctx,
			concreteTypes,
			target,
			cont,
		),
	)
}

func unmarshalTuple(
	ctx Ctx,
	concreteTypes []reflect.Type,
	target reflect.Value,
	cont Sink,
) Sink {

	var values []reflect.Value
	var valueTypes []reflect.Type
	var sink Sink
	sink = func(token *Token) (Sink, error) {
		if token == nil {
			return nil, UnmarshalError{ExpectingValue}
		}

		if token.Kind == KindTupleEnd {

			// too few values
			if len(concreteTypes) > 0 {
				return nil, UnmarshalError{ExpectingValue}
			}

			targetType := target.Type()
			if targetType.Kind() == reflect.Func {
				// arg nums not match
				if !targetType.IsVariadic() && targetType.NumIn() != len(values) {
					return nil, UnmarshalError{BadTupleType}
				}
				if !target.IsNil() {
					rets := target.Call(values)
					for _, ret := range rets {
						if e, ok := ret.Interface().(error); ok {
							return nil, e
						}
					}
				}

			} else {
				// not func type, set func() (...) tuple
				if len(values) > 50 {
					return nil, UnmarshalError{TooManyElement}
				}
				funcType := reflect.FuncOf(
					[]reflect.Type{},
					valueTypes,
					false,
				)
				if !funcType.AssignableTo(target.Elem().Type()) {
					return nil, UnmarshalError{BadTupleType}
				}
				target.Elem().Set(reflect.MakeFunc(
					funcType,
					func(args []reflect.Value) []reflect.Value {
						return values
					},
				))
			}

			return cont, nil
		}

		// collect values
		if len(concreteTypes) > 0 {
			t := concreteTypes[0]
			concreteTypes = concreteTypes[1:]
			value := reflect.New(t)
			valueTypes = append(valueTypes, t)
			return ctx.Unmarshal(
				ctx,
				value,
				func(token *Token) (Sink, error) {
					values = append(values, value.Elem())
					return sink(token)
				},
			)(token)

		} else {
			var obj any
			value := reflect.ValueOf(&obj)
			return ctx.Unmarshal(
				ctx,
				value,
				func(token *Token) (Sink, error) {
					if obj != nil {
						values = append(values, value.Elem().Elem())
						valueTypes = append(valueTypes, reflect.TypeOf(obj))
					} else {
						values = append(values, value.Elem())
						valueTypes = append(valueTypes, value.Type().Elem())
					}
					return sink(token)
				},
			)(token)

		}

	}

	return sink
}
