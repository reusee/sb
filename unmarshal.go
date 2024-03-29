package sb

import (
	"encoding"
	"fmt"
	gotoken "go/token"
	"io"
	"math"
	"reflect"
	"strconv"

	"github.com/reusee/e5"
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

func TapUnmarshal(ctx Ctx, target any, fn func(Ctx, Token, reflect.Value)) Sink {
	unmarshal := func(ctx Ctx, target reflect.Value, cont Sink) Sink {
		return func(token *Token) (Sink, error) {
			if token.Invalid() {
				return cont, nil
			}
			fn(ctx, *token, target)
			return UnmarshalValue(ctx, target, cont)(token)
		}
	}
	ctx.Unmarshal = unmarshal
	return unmarshal(ctx, reflect.ValueOf(target), nil)
}

type UnmarshalFunc func(Ctx, Sink) Sink

var _ SBUnmarshaler = UnmarshalFunc(nil)

func (f UnmarshalFunc) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return f(ctx, cont)
}

func UnmarshalValue(ctx Ctx, target reflect.Value, cont Sink) Sink {
	if ctx.Unmarshal == nil {
		ctx.Unmarshal = UnmarshalValue
	}

	return func(token *Token) (next Sink, err error) {
		defer func() {
			err = we.With(WithPath(ctx))(err)
		}()

		// convert literal token
		if token.Valid() && token.Kind == KindLiteral {
			copy := *token
			token = &copy

			if !target.IsValid() {
				token.Kind = KindString

			} else {
				targetType := target.Type()
				if targetType.Kind() != reflect.Ptr {
					return nil, we.With(BadTargetType)(UnmarshalError)
				}
				valueKind := targetType.Elem().Kind()
				switch valueKind {

				case reflect.String:
					token.Kind = KindString

				case reflect.Bool:
					token.Kind = KindBool
					b, err := strconv.ParseBool(token.Value.(string))
					if err != nil {
						return nil, we.With(UnmarshalError)(err)
					}
					token.Value = b

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					i, err := strconv.ParseInt(token.Value.(string), 10, 64)
					if err != nil {
						return nil, we.With(UnmarshalError)(err)
					}
					switch valueKind {
					case reflect.Int:
						token.Kind = KindInt
						token.Value = int(i)
					case reflect.Int8:
						token.Kind = KindInt8
						token.Value = int8(i)
					case reflect.Int16:
						token.Kind = KindInt16
						token.Value = int16(i)
					case reflect.Int32:
						token.Kind = KindInt32
						token.Value = int32(i)
					case reflect.Int64:
						token.Kind = KindInt64
						token.Value = int64(i)
					}

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					u, err := strconv.ParseUint(token.Value.(string), 10, 64)
					if err != nil {
						return nil, we.With(UnmarshalError)(err)
					}
					switch valueKind {
					case reflect.Uint:
						token.Kind = KindUint
						token.Value = uint(u)
					case reflect.Uint8:
						token.Kind = KindUint8
						token.Value = uint8(u)
					case reflect.Uint16:
						token.Kind = KindUint16
						token.Value = uint16(u)
					case reflect.Uint32:
						token.Kind = KindUint32
						token.Value = uint32(u)
					case reflect.Uint64:
						token.Kind = KindUint64
						token.Value = uint64(u)
					case reflect.Uintptr:
						token.Kind = KindPointer
						token.Value = uintptr(u)
					}

				case reflect.Float32:
					f, err := strconv.ParseFloat(token.Value.(string), 32)
					if err != nil {
						return nil, we.With(UnmarshalError)(err)
					}
					token.Kind = KindFloat32
					token.Value = float32(f)

				case reflect.Float64:
					f, err := strconv.ParseFloat(token.Value.(string), 64)
					if err != nil {
						return nil, we.With(UnmarshalError)(err)
					}
					token.Kind = KindFloat64
					token.Value = f

				default:
					return nil, we.With(BadTargetType)(UnmarshalError)

				}
			}

		}

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

			case *SBUnmarshaler:
				if v != nil {
					return (*v).UnmarshalSB(ctx, cont)(token)
				}

			case SBUnmarshaler:
				return v.UnmarshalSB(ctx, cont)(token)

			case encoding.BinaryUnmarshaler:
				if token.Invalid() {
					return nil, we.With(
						TypeMismatch(KindInvalid, reflect.String),
						io.ErrUnexpectedEOF,
					)(
						UnmarshalError,
					)
				}
				if token.Kind != KindString {
					return nil, we.With(TypeMismatch(token.Kind, reflect.String))(UnmarshalError)
				}
				if err := v.UnmarshalBinary(
					[]byte(token.Value.(string)),
				); err != nil {
					return nil, we.With(UnmarshalError)(err)
				}
				return cont, nil

			case encoding.TextUnmarshaler:
				if token.Invalid() {
					return nil, we.With(
						TypeMismatch(KindInvalid, reflect.String),
						io.ErrUnexpectedEOF,
					)(
						UnmarshalError,
					)
				}
				if token.Kind != KindString {
					return nil, we.With(TypeMismatch(token.Kind, reflect.String))(UnmarshalError)
				}
				if err := v.UnmarshalText(
					[]byte(token.Value.(string)),
				); err != nil {
					return nil, we.With(UnmarshalError)(err)
				}
				return cont, nil

			}
		}

		if token.Invalid() {
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(
				UnmarshalError,
			)
		}

		switch token.Kind {
		case KindNil:
			return cont, nil
		case KindArrayEnd, KindObjectEnd, KindMapEnd, KindTupleEnd:
			return nil, we.With(UnexpectedEndToken)(UnmarshalError)
		}

		targetType := target.Type()
		targetKind := targetType.Kind()
		var valueType reflect.Type
		var valueKind reflect.Kind
		switch targetKind {
		case reflect.Func:
			valueType = targetType
			valueKind = targetKind
		case reflect.Ptr:
			valueType = targetType.Elem()
			valueKind = valueType.Kind()
		default:
			return nil, we.With(BadTargetType)(UnmarshalError)
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
					return cont.Sink(token)
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
					return nil, we.With(TypeMismatch(KindBool, valueKind))(UnmarshalError)
				}
				target.Elem().SetBool(token.Value.(bool))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(bool)))
			}

		case KindInt:
			if hasConcreteType {
				if valueKind != reflect.Int {
					return nil, we.With(TypeMismatch(KindInt, valueKind))(UnmarshalError)
				}
				target.Elem().SetInt(int64(token.Value.(int)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int)))
			}

		case KindInt8:
			if hasConcreteType {
				if valueKind != reflect.Int8 {
					return nil, we.With(TypeMismatch(KindInt8, valueKind))(UnmarshalError)
				}
				target.Elem().SetInt(int64(token.Value.(int8)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int8)))
			}

		case KindInt16:
			if hasConcreteType {
				if valueKind != reflect.Int16 {
					return nil, we.With(TypeMismatch(KindInt16, valueKind))(UnmarshalError)
				}
				target.Elem().SetInt(int64(token.Value.(int16)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int16)))
			}

		case KindInt32:
			if hasConcreteType {
				if valueKind != reflect.Int32 {
					return nil, we.With(TypeMismatch(KindInt32, valueKind))(UnmarshalError)
				}
				target.Elem().SetInt(int64(token.Value.(int32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int32)))
			}

		case KindInt64:
			if hasConcreteType {
				if valueKind != reflect.Int64 {
					return nil, we.With(TypeMismatch(KindInt64, valueKind))(UnmarshalError)
				}
				target.Elem().SetInt(token.Value.(int64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(int64)))
			}

		case KindUint:
			if hasConcreteType {
				if valueKind != reflect.Uint {
					return nil, we.With(TypeMismatch(KindUint, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(uint64(token.Value.(uint)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint)))
			}

		case KindUint8:
			if hasConcreteType {
				if valueKind != reflect.Uint8 {
					return nil, we.With(TypeMismatch(KindUint8, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(uint64(token.Value.(uint8)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint8)))
			}

		case KindUint16:
			if hasConcreteType {
				if valueKind != reflect.Uint16 {
					return nil, we.With(TypeMismatch(KindUint16, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(uint64(token.Value.(uint16)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint16)))
			}

		case KindUint32:
			if hasConcreteType {
				if valueKind != reflect.Uint32 {
					return nil, we.With(TypeMismatch(KindUint32, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(uint64(token.Value.(uint32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint32)))
			}

		case KindUint64:
			if hasConcreteType {
				if valueKind != reflect.Uint64 {
					return nil, we.With(TypeMismatch(KindUint64, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(token.Value.(uint64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uint64)))
			}

		case KindPointer:
			if hasConcreteType {
				if valueKind != reflect.Uintptr {
					return nil, we.With(TypeMismatch(KindPointer, valueKind))(UnmarshalError)
				}
				target.Elem().SetUint(uint64(token.Value.(uintptr)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(uintptr)))
			}

		case KindFloat32:
			if hasConcreteType {
				if valueKind != reflect.Float32 {
					return nil, we.With(TypeMismatch(KindFloat32, valueKind))(UnmarshalError)
				}
				target.Elem().SetFloat(float64(token.Value.(float32)))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(float32)))
			}

		case KindFloat64:
			if hasConcreteType {
				if valueKind != reflect.Float64 {
					return nil, we.With(TypeMismatch(KindFloat64, valueKind))(UnmarshalError)
				}
				target.Elem().SetFloat(token.Value.(float64))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(float64)))
			}

		case KindNaN:
			if hasConcreteType {
				if valueKind != reflect.Float32 && valueKind != reflect.Float64 {
					return nil, we.With(TypeMismatch(KindNaN, valueKind))(UnmarshalError)
				}
				target.Elem().SetFloat(math.NaN())
			} else {
				target.Elem().Set(reflect.ValueOf(math.NaN()))
			}

		case KindString:
			if hasConcreteType {
				if valueKind != reflect.String {
					return nil, we.With(TypeMismatch(KindString, valueKind))(UnmarshalError)
				}
				target.Elem().SetString(token.Value.(string))
			} else {
				target.Elem().Set(reflect.ValueOf(token.Value.(string)))
			}

		case KindBytes:
			if hasConcreteType {
				if !isBytes(valueType) {
					return nil, we.With(TypeMismatch(KindBytes, valueKind))(UnmarshalError)
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

				switch valueKind {
				case reflect.Array:
					// array
					return UnmarshalArray(
						ctx,
						target,
						cont,
					)(token)

				case reflect.Slice:
					// slice
					return UnmarshalSlice(
						ctx,
						target,
						valueType,
						cont,
					)(token)

				default:
					return nil, we.With(TypeMismatch(KindArray, valueKind))(UnmarshalError)
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
					return nil, we.With(TypeMismatch(KindObject, valueKind))(UnmarshalError)
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
					return nil, we.With(TypeMismatch(KindMap, valueKind))(UnmarshalError)
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
					return nil, we.With(TypeMismatch(KindTuple, valueKind))(UnmarshalError)
				}
			}
			return UnmarshalTuple(
				ctx,
				target,
				valueType,
				cont,
			)(token)

		case KindTypeName:
			if hasConcreteType {
				return notNull(ctx, ctx.Unmarshal(ctx, target, cont)), nil
			}
			t, ok := registeredNameToType.Load(token.Value.(string))
			if !ok {
				return notNull(ctx, ctx.Unmarshal(ctx, target, cont)), nil
			}
			v := reflect.New(t.(reflect.Type))
			return notNull(ctx, ctx.Unmarshal(
				ctx,
				v,
				func(token *Token) (Sink, error) {
					target.Elem().Set(v.Elem())
					return cont.Sink(token)
				},
			)), nil

		default:
			return nil, we.With(BadTokenKind, token.Kind)(UnmarshalError)
		}

		return cont, nil
	}

}

func notNull(ctx Ctx, cont Sink) Sink {
	return func(p *Token) (Sink, error) {
		if p == nil {
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}
		return cont(p)
	}
}

func ExpectKind(ctx Ctx, kind Kind, cont Sink) Sink {
	return func(token *Token) (next Sink, err error) {
		defer func() {
			err = we.With(WithPath(ctx))(err)
		}()
		if token.Invalid() {
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		} else if token.Kind != kind {
			return nil, we.With(
				e5.Info("expecting %s, got %s", kind, token.Kind),
			)(UnmarshalError)
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}
		if p.Kind == KindArrayEnd {
			return cont, nil
		}
		if idx >= target.Elem().Len() {
			return nil, we.With(WithPath(ctx), TooManyElement)(UnmarshalError)
		}

		e := target.Elem().Index(idx).Addr()
		idx++
		return ctx.Unmarshal(
			ctx.WithPath(idx-1),
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}
		if p.Kind == KindArrayEnd {
			target.Elem().Set(slice)
			return cont, nil
		}
		elemPtr := reflect.New(valueType.Elem())
		slice = reflect.Append(slice, elemPtr.Elem())

		return ctx.Unmarshal(
			ctx.WithPath(slice.Len()-1),
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}
		if p.Kind == KindArrayEnd {
			target.Elem().Set(reflect.ValueOf(slice))
			return cont, nil
		}

		var value any
		return ctx.Unmarshal(
			ctx.WithPath(len(slice)),
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}
		if p.Kind == KindObjectEnd {
			return cont, nil
		}
		var name string

		return ctx.Unmarshal(
			ctx,
			reflect.ValueOf(&name),
			func(token *Token) (Sink, error) {
				field, ok := valueType.FieldByNameFunc(func(str string) bool {
					return str == name
				})
				if !ok {
					if ctx.DisallowUnknownStructFields {
						// check field deprecation
						if fieldIsDeprecated(valueType, name) {
							// skip next value
							var value any
							return ctx.Unmarshal(
								ctx.WithPath(name),
								reflect.ValueOf(&value),
								sink,
							)(token)
						}
						return nil, we.With(
							WithPath(ctx),
							UnknownFieldName,
							fmt.Errorf("field: %s", name),
						)(UnmarshalError)
					} else {
						// skip next value
						var value any
						return ctx.Unmarshal(
							ctx.WithPath(name),
							reflect.ValueOf(&value),
							sink,
						)(token)
					}

				} else {
					return ctx.Unmarshal(
						ctx.WithPath(target.Elem().Type().FieldByIndex(field.Index).Name),
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
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
					return nil, we.With(WithPath(ctx), BadFieldName, fmt.Errorf("field: %s", name))(UnmarshalError)
				}
				if _, ok := names[name]; ok {
					return nil, we.With(WithPath(ctx), DuplicatedFieldName)(UnmarshalError)
				}
				names[name] = struct{}{}
				var value any

				return ctx.Unmarshal(
					ctx.WithPath(name),
					reflect.ValueOf(&value),
					func(token *Token) (Sink, error) {
						if value == nil {
							return nil, we.With(
								WithPath(ctx),
								io.ErrUnexpectedEOF,
							)(UnmarshalError)
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
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
					ctx.WithPath(key.Elem().Interface()),
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
		ctx,
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
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
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
				key = toComparable(key)
				if key == nil {
					return nil, we.With(WithPath(ctx), BadMapKey)(UnmarshalError)
				} else if !reflect.TypeOf(key).Comparable() {
					return nil, we.With(WithPath(ctx), BadMapKey)(UnmarshalError)
				} else if f, ok := key.(float64); ok && f != f {
					return nil, we.With(WithPath(ctx), BadMapKey)(UnmarshalError)
				} else if f, ok := key.(float32); ok && f != f {
					return nil, we.With(WithPath(ctx), BadMapKey)(UnmarshalError)
				}
				var value any

				return ctx.Unmarshal(
					ctx.WithPath(key),
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

func toComparable(value any) any {
	if slice, ok := value.([]byte); ok {
		array := reflect.New(
			reflect.ArrayOf(
				len(slice),
				reflect.TypeOf((*byte)(nil)).Elem(),
			),
		).Elem()
		reflect.Copy(
			array.Slice(0, len(slice)),
			reflect.ValueOf(slice),
		)
		return array.Interface()
	}
	return value
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
		ctx,
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
		if token.Invalid() {
			return nil, we.With(
				WithPath(ctx),
				io.ErrUnexpectedEOF,
			)(UnmarshalError)
		}

		if token.Kind == KindTupleEnd {

			// too few values
			if len(concreteTypes) > 0 {
				return nil, we.With(WithPath(ctx), TooFewElement)(UnmarshalError)
			}

			targetType := target.Type()
			if targetType.Kind() == reflect.Func {
				// arg nums not match
				if !targetType.IsVariadic() && targetType.NumIn() != len(values) {
					return nil, we.With(WithPath(ctx), BadTupleType)(UnmarshalError)
				}
				if !target.IsNil() {
					rets := target.Call(values)
					for _, ret := range rets {
						if e, ok := ret.Interface().(error); ok {
							return nil, we.With(UnmarshalError, WithPath(ctx))(e)
						}
					}
				}

			} else {
				// not func type, set func() (...) tuple
				if len(values) > 50 {
					return nil, we.With(WithPath(ctx), TooManyElement)(UnmarshalError)
				}
				funcType := reflect.FuncOf(
					[]reflect.Type{},
					valueTypes,
					false,
				)
				if !funcType.AssignableTo(target.Elem().Type()) {
					return nil, we.With(WithPath(ctx), BadTupleType)(UnmarshalError)
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
				ctx.WithPath(len(valueTypes)-1),
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
				ctx.WithPath(len(valueTypes)),
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
