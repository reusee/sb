package sb

import (
	"encoding"
	"math"
	"reflect"
	"sort"
	"sync"
)

type SBMarshaler interface {
	MarshalSB(vm ValueMarshalFunc, cont Proc) Proc
}

func Marshal(value any) *Proc {
	marshaler := MarshalValue(MarshalValue, reflect.ValueOf(value), nil)
	return &marshaler
}

type ValueMarshalFunc func(
	fn ValueMarshalFunc,
	value reflect.Value,
	cont Proc,
) (
	proc Proc,
)

func MarshalAny(vm ValueMarshalFunc, value any, cont Proc) Proc {
	return vm(vm, reflect.ValueOf(value), cont)
}

func MarshalValue(vm ValueMarshalFunc, value reflect.Value, cont Proc) Proc {
	return func() (*Token, Proc, error) {

		if value.IsValid() {
			i := value.Interface()
			if v, ok := i.(SBMarshaler); ok {
				return nil, v.MarshalSB(vm, cont), nil
			} else if v, ok := i.(encoding.BinaryMarshaler); ok {
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, nil, err
				}
				return nil, vm(vm, reflect.ValueOf(string(bs)), cont), nil
			} else if v, ok := i.(encoding.TextMarshaler); ok {
				bs, err := v.MarshalText()
				if err != nil {
					return nil, nil, err
				}
				return nil, vm(vm, reflect.ValueOf(string(bs)), cont), nil
			} else if v, ok := i.(*Token); ok {
				return v, cont, nil
			}
		}

		switch value.Kind() {

		case reflect.Invalid:
			return &Token{
				Kind: KindNil,
			}, cont, nil

		case reflect.Ptr, reflect.Interface:
			if value.IsNil() {
				return &Token{
					Kind: KindNil,
				}, cont, nil
			} else {
				return nil, vm(vm, value.Elem(), cont), nil
			}

		case reflect.Bool:
			return &Token{
				Kind:  KindBool,
				Value: bool(value.Bool()),
			}, cont, nil

		case reflect.Int:
			return &Token{
				Kind:  KindInt,
				Value: int(value.Int()),
			}, cont, nil

		case reflect.Int8:
			return &Token{
				Kind:  KindInt8,
				Value: int8(value.Int()),
			}, cont, nil

		case reflect.Int16:
			return &Token{
				Kind:  KindInt16,
				Value: int16(value.Int()),
			}, cont, nil

		case reflect.Int32:
			return &Token{
				Kind:  KindInt32,
				Value: int32(value.Int()),
			}, cont, nil

		case reflect.Int64:
			return &Token{
				Kind:  KindInt64,
				Value: int64(value.Int()),
			}, cont, nil

		case reflect.Uint:
			return &Token{
				Kind:  KindUint,
				Value: uint(value.Uint()),
			}, cont, nil

		case reflect.Uint8:
			return &Token{
				Kind:  KindUint8,
				Value: uint8(value.Uint()),
			}, cont, nil

		case reflect.Uint16:
			return &Token{
				Kind:  KindUint16,
				Value: uint16(value.Uint()),
			}, cont, nil

		case reflect.Uint32:
			return &Token{
				Kind:  KindUint32,
				Value: uint32(value.Uint()),
			}, cont, nil

		case reflect.Uint64:
			return &Token{
				Kind:  KindUint64,
				Value: uint64(value.Uint()),
			}, cont, nil

		case reflect.Float32:
			if math.IsNaN(value.Float()) {
				return &Token{
					Kind: KindNaN,
				}, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat32,
					Value: float32(value.Float()),
				}, cont, nil
			}

		case reflect.Float64:
			if math.IsNaN(value.Float()) {
				return &Token{
					Kind: KindNaN,
				}, cont, nil
			} else {
				return &Token{
					Kind:  KindFloat64,
					Value: float64(value.Float()),
				}, cont, nil
			}

		case reflect.Array, reflect.Slice:
			if isBytes(value.Type()) {
				return &Token{
					Kind:  KindBytes,
					Value: toBytes(value),
				}, cont, nil
			} else {
				return &Token{
					Kind: KindArray,
				}, MarshalArray(vm, value, 0, cont), nil
			}

		case reflect.String:
			return &Token{
				Kind:  KindString,
				Value: value.String(),
			}, cont, nil

		case reflect.Struct:
			return &Token{
				Kind: KindObject,
			}, MarshalStruct(vm, value, cont), nil

		case reflect.Map:
			return &Token{
				Kind: KindMap,
			}, MarshalMap(vm, value, cont), nil

		case reflect.Func:
			items := value.Call([]reflect.Value{})
			return &Token{
					Kind: KindTuple,
				}, MarshalTuple(
					vm,
					items,
					cont,
				), nil

		default:
			return nil, cont, nil

		}
	}
}

var arrayEndToken = reflect.ValueOf(&Token{
	Kind: KindArrayEnd,
})

func MarshalArray(vm ValueMarshalFunc, value reflect.Value, index int, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if index >= value.Len() {
			return nil, vm(
				vm,
				arrayEndToken,
				cont,
			), nil
		}
		v := value.Index(index)
		index++
		return nil, vm(
			vm,
			v,
			proc,
		), nil
	}
	return proc
}

var structFields sync.Map

func MarshalStruct(vm ValueMarshalFunc, value reflect.Value, cont Proc) Proc {
	var fields []reflect.StructField
	valueType := value.Type()
	if v, ok := structFields.Load(valueType); ok {
		fields = v.([]reflect.StructField)
	} else {
		numField := valueType.NumField()
		for i := 0; i < numField; i++ {
			field := valueType.Field(i)
			if field.PkgPath == "" {
				fields = append(fields, field)
			}
		}
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})
		structFields.Store(valueType, fields)
	}
	return MarshalStructFields(vm, value, fields, cont)
}

func MarshalStructNonEmpty(vm ValueMarshalFunc, value reflect.Value, cont Proc) Proc {
	var fields []reflect.StructField
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := value.Field(i)
		if fieldValue.IsZero() {
			continue
		}
		field := t.Field(i)
		if field.Type.Kind() == reflect.Slice && fieldValue.Len() == 0 {
			continue
		}
		fields = append(fields, field)
	}
	return func() (*Token, Proc, error) {
		return &Token{
			Kind: KindObject,
		}, MarshalStructFields(vm, value, fields, cont), nil
	}
}

var objectEndToken = reflect.ValueOf(&Token{
	Kind: KindObjectEnd,
})

func MarshalStructFields(vm ValueMarshalFunc, value reflect.Value, fields []reflect.StructField, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(fields) == 0 {
			return nil, vm(
				vm,
				objectEndToken,
				cont,
			), nil
		}
		field := fields[0]
		fields = fields[1:]
		return nil, vm(
			vm,
			reflect.ValueOf(field.Name),
			func() (*Token, Proc, error) {
				return nil, vm(
					vm,
					value.FieldByIndex(field.Index),
					proc,
				), nil
			},
		), nil
	}
	return proc
}

type MapTuple struct {
	KeyTokens Tokens
	Key       reflect.Value
	Value     reflect.Value
}

func MarshalMap(vm ValueMarshalFunc, value reflect.Value, cont Proc) Proc {
	return MarshalMapIter(
		vm,
		value,
		value.MapRange(),
		make([]*MapTuple, 0, value.Len()),
		cont,
	)
}

func MarshalMapIter(vm ValueMarshalFunc, value reflect.Value, iter *reflect.MapIter, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if !iter.Next() {
			// done
			sort.Slice(tuples, func(i, j int) bool {
				return MustCompare(
					tuples[i].KeyTokens.Iter(),
					tuples[j].KeyTokens.Iter(),
				) < 0
			})
			return nil, MarshalMapTuples(vm, tuples, cont), nil
		}
		var tokens Tokens
		if err := Copy(
			Marshal(iter.Key().Interface()),
			CollectTokens(&tokens),
		); err != nil {
			return nil, nil, err
		}
		if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, nil, MarshalError{BadMapKey}
		}
		tuples = append(tuples, &MapTuple{
			KeyTokens: tokens,
			Key:       iter.Key(),
			Value:     iter.Value(),
		})
		return nil, proc, nil
	}
	return proc
}

var mapEndToken = reflect.ValueOf(&Token{
	Kind: KindMapEnd,
})

func MarshalMapTuples(vm ValueMarshalFunc, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(tuples) == 0 {
			return nil, vm(
				vm,
				mapEndToken,
				cont,
			), nil
		}
		tuple := tuples[0]
		tuples = tuples[1:]
		return nil, vm(
			vm,
			tuple.Key,
			func() (*Token, Proc, error) {
				return nil, vm(
					vm,
					tuple.Value,
					proc,
				), nil
			},
		), nil
	}
	return proc
}

var tupleEndToken = reflect.ValueOf(&Token{
	Kind: KindTupleEnd,
})

func MarshalTuple(vm ValueMarshalFunc, items []reflect.Value, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(items) == 0 {
			return nil, vm(
				vm,
				tupleEndToken,
				cont,
			), nil
		} else {
			v := items[0]
			items = items[1:]
			return nil, vm(
				vm,
				v,
				proc,
			), nil
		}
	}
	return proc
}
