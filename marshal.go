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

func MarshalTokens(tokens []Token, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(tokens) == 0 {
			return nil, cont, nil
		}
		v := tokens[0]
		tokens = tokens[1:]
		return &v, proc, nil
	}
	return proc
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

func MarshalArray(vm ValueMarshalFunc, value reflect.Value, index int, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if index >= value.Len() {
			return &Token{
				Kind: KindArrayEnd,
			}, cont, nil
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
			fields = append(fields, valueType.Field(i))
		}
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})
		structFields.Store(valueType, fields)
	}
	return MarshalStructFields(vm, value, fields, cont)
}

func MarshalStructFields(vm ValueMarshalFunc, value reflect.Value, fields []reflect.StructField, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(fields) == 0 {
			return &Token{
				Kind: KindObjectEnd,
			}, cont, nil
		}
		field := fields[0]
		fields = fields[1:]
		return &Token{
				Kind:  KindString,
				Value: field.Name,
			}, vm(
				vm,
				value.FieldByIndex(field.Index),
				proc,
			), nil
	}
	return proc
}

type MapTuple struct {
	KeyTokens Tokens
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
			return nil, MarshalMapValue(vm, tuples, cont), nil
		}
		marshaler := vm(vm, iter.Key(), nil)
		tokens, err := TokensFromStream(&marshaler)
		if err != nil {
			return nil, nil, err
		} else if len(tokens) == 0 ||
			(len(tokens) == 1 && tokens[0].Kind == KindNaN) {
			return nil, nil, MarshalError{BadMapKey}
		}
		tuples = append(tuples, &MapTuple{
			KeyTokens: tokens,
			Value:     iter.Value(),
		})
		return nil, proc, nil
	}
	return proc
}

func MarshalMapValue(vm ValueMarshalFunc, tuples []*MapTuple, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(tuples) == 0 {
			return &Token{
				Kind: KindMapEnd,
			}, cont, nil
		}
		tuple := tuples[0]
		tuples = tuples[1:]
		return nil, MarshalTokens(
			tuple.KeyTokens,
			vm(
				vm,
				tuple.Value,
				proc,
			),
		), nil
	}
	return proc
}

func MarshalTuple(vm ValueMarshalFunc, items []reflect.Value, cont Proc) Proc {
	var proc Proc
	proc = func() (*Token, Proc, error) {
		if len(items) == 0 {
			return &Token{
				Kind: KindTupleEnd,
			}, cont, nil
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
