package sb

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

type UnmarshalTestCase struct {
	value  any
	target any
	err    error
}

type testBytes []byte

var unmarshalTestCases = []UnmarshalTestCase{
	0:  {true, bool(false), nil},
	1:  {int8(42), int8(0), nil},
	2:  {int16(42), int16(0), nil},
	3:  {int32(42), int32(0), nil},
	4:  {int64(42), int64(0), nil},
	5:  {uint(42), uint(0), nil},
	6:  {uint8(42), uint8(0), nil},
	7:  {uint16(42), uint16(0), nil},
	8:  {uint32(42), uint32(0), nil},
	9:  {uint64(42), uint64(0), nil},
	10: {float32(42), float32(0), nil},
	11: {float64(42), float64(0), nil},
	12: {string("42"), string(""), nil},
	13: {map[int]int{1: 1}, map[int]int(nil), nil},
	14: {[]byte("foo"), []byte("foo"), nil},
	15: {testBytes{42}, testBytes{42}, nil},
	16: {[3]int{1}, [3]int{1}, nil},
	17: {
		func() (int, string) {
			return 42, "42"
		},
		(func() (int, string))(nil),
		nil,
	},
	18: {true, int(0), ExpectingBool},
	19: {42, true, ExpectingInt},
	20: {int8(42), true, ExpectingInt8},
	21: {int16(42), true, ExpectingInt16},
	22: {int32(42), true, ExpectingInt32},
	23: {int64(42), true, ExpectingInt64},
	24: {uint(42), true, ExpectingUint},
	25: {uint8(42), true, ExpectingUint8},
	26: {uint16(42), true, ExpectingUint16},
	27: {uint32(42), true, ExpectingUint32},
	28: {uint64(42), true, ExpectingUint64},
	29: {float32(42), true, ExpectingFloat32},
	30: {float64(42), true, ExpectingFloat64},
	31: {math.NaN(), true, ExpectingFloat},
	32: {"42", true, ExpectingString},
	33: {[]int{42}, true, ExpectingSequence},
	34: {map[int]int{}, true, ExpectingMap},
	35: {42, (****string)(nil), ExpectingInt},
	36: {[]byte("foo"), true, ExpectingBytes},
	37: {
		func() int { return 42 },
		true,
		ExpectingTuple,
	},
	38: {
		func() int { return 42 },
		(func(int) int)(nil),
		BadTupleType,
	},
	39: {testBool(true), testBool(false), nil},
	40: {testUint8(42), testUint8(0), nil},
	41: {testInt8(42), testInt8(0), nil},
	42: {testInt16(42), testInt16(0), nil},
	43: {testInt32(42), testInt32(0), nil},
	44: {testInt64(42), testInt64(0), nil},
	45: {testUint(42), testUint(0), nil},
	46: {testUint16(42), testUint16(0), nil},
	47: {testUint64(42), testUint64(0), nil},
	48: {testFloat32(42), testFloat32(0), nil},
	49: {testFloat64(42), testFloat64(0), nil},
	50: {testString("foo"), testString(""), nil},
}

func TestUnmarshal(t *testing.T) {
	for i, c := range unmarshalTestCases {
		stream := Marshal(c.value)
		ptr := reflect.New(reflect.TypeOf(c.target))
		err := Copy(stream, Unmarshal(ptr.Interface()))
		if !errors.Is(err, c.err) {
			t.Fatalf("case %d, expecting %v, got %v", i, c.err, err)
		}
		if err == nil {
			if ptr.Elem().Kind() == reflect.Func {
				var items1 []any
				for _, v := range reflect.ValueOf(c.value).Call([]reflect.Value{}) {
					items1 = append(items1, v.Interface())
				}
				var items2 []any
				for _, v := range ptr.Elem().Call([]reflect.Value{}) {
					items2 = append(items2, v.Interface())
				}
				if !reflect.DeepEqual(items1, items2) {
					t.Fatal()
				}
			} else {
				if !reflect.DeepEqual(c.value, ptr.Elem().Interface()) {
					t.Fatalf("not equal: %#v, %#v", c.value, ptr.Elem().Interface())
				}
			}
		}
	}
}

func TestUnmarshalNaN(t *testing.T) {
	stream := Marshal(math.NaN())
	var f float64
	if err := Copy(stream, Unmarshal(&f)); err != nil {
		t.Fatal(err)
	}
	if f == f {
		t.Fatal()
	}
}

func TestUnmarshalArray(t *testing.T) {
	type foo [3]byte
	type S struct {
		Foos []foo
	}
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(S{
			Foos: []foo{
				foo{1},
				foo{2},
			},
		}),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}
	var s S
	if err := Copy(Decode(buf), Unmarshal(&s)); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalNamedUint(t *testing.T) {
	type Foo uint32
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(Foo(42)),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}
	var foo Foo
	if err := Copy(Decode(buf), Unmarshal(&foo)); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalStructWithPrivateField(t *testing.T) {
	type Foo struct {
		Bar int
		Foo int
	}
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(Foo{42, 42}),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}
	type Bar struct {
		bar int
		Foo int
	}
	var bar Bar
	if err := Copy(Decode(buf), Unmarshal(&bar)); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalIncompleteStream(t *testing.T) {
	cases := []Tokens{
		{
			{Kind: KindObject},
		},
		{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
		},
		{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindString, Value: "Bar"},
			{KindInt, 42},
		},
		{},
		{
			{Kind: KindArray},
		},
		{
			{Kind: KindMap},
		},
		{
			{Kind: KindTuple},
		},
	}

	for i, c := range cases {
		var v any
		err := Copy(c.Iter(), Unmarshal(&v))
		if err == nil {
			t.Fatalf("shoud error: %d", i)
		}
	}

}

type badBinaryUnmarshaler struct{}

var _ encoding.BinaryUnmarshaler = new(badBinaryUnmarshaler)

func (b *badBinaryUnmarshaler) UnmarshalBinary(data []byte) error {
	return fmt.Errorf("bad")
}

func TestUnmarshalBadBinaryUnmarshaler(t *testing.T) {
	var b badBinaryUnmarshaler

	// bad decoder
	if err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindString),
		})),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// no token
	if err := Copy(
		Decode(bytes.NewReader(nil)),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// bad token
	if err := Copy(
		Tokens{
			{KindInt, 42},
		}.Iter(),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// bad unmarshaler
	if err := Copy(
		Tokens{
			{Kind: KindString, Value: "foo"},
		}.Iter(),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

}

type badTextUnmarshaler struct{}

var _ encoding.TextUnmarshaler = new(badTextUnmarshaler)

func (b *badTextUnmarshaler) UnmarshalText(data []byte) error {
	return fmt.Errorf("bad")
}

func TestUnmarshalBadTextUnmarshaler(t *testing.T) {
	var b badTextUnmarshaler

	// bad decoder
	if err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindString),
		})),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// no token
	if err := Copy(
		Decode(bytes.NewReader(nil)),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// bad token
	if err := Copy(
		Tokens{
			{KindInt, 42},
		}.Iter(),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

	// bad unmarshaler
	if err := Copy(
		Tokens{
			{Kind: KindString, Value: "foo"},
		}.Iter(),
		Unmarshal(&b),
	); err == nil {
		t.Fatal()
	}

}

func TestUnmarshalToNilPtr(t *testing.T) {
	if err := Copy(
		Tokens{
			{Kind: KindInt, Value: 42},
		}.Iter(),
		Unmarshal((*int)(nil)),
	); err != nil {
		t.Fatal(err)
	}
}

func TestBadArray(t *testing.T) {
	var v [2]int

	// bad decoder
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindArray),
			byte(KindString),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindArray),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// too many element
	err = Copy(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{KindInt, 42},
			{KindInt, 42},
			{KindInt, 42},
		}.Iter(),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{Kind: KindInt8, Value: int8(42)},
		}.Iter(),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}
}

func TestBadSlice(t *testing.T) {
	var v []int

	// bad decoder
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindArray),
			byte(KindString),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindArray),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{Kind: KindInt8, Value: int8(42)},
		}.Iter(),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}
}

func TestBadObject(t *testing.T) {
	var v struct {
		Foo int
	}

	// bad decoder
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindObject),
			byte(KindString),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindObject),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindObject},
			{KindInt, 42},
		}.Iter(),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindObject},
			{Kind: KindString, Value: "42"},
			{Kind: KindInt8, Value: int8(42)},
		}.Iter(),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindObject),
			byte(KindString), 0,
			byte(KindString),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindObject),
			byte(KindString), 3, 'F', 'o', 'o',
			byte(KindString),
		})),
		Unmarshal(&v),
	)
	if err == nil {
		t.Fatal(err)
	}

}

func TestBadMap(t *testing.T) {
	var m map[int]int

	// bad decoder
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString),
		})),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
		})),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindMap},
			{Kind: KindString, Value: "foo"},
		}.Iter(),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Copy(
		Tokens{
			{Kind: KindMap},
			{Kind: KindString, Value: "42"},
			{Kind: KindInt8, Value: int8(42)},
		}.Iter(),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString), 0,
			byte(KindString),
		})),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString), 3, 'F', 'o', 'o',
			byte(KindString),
		})), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindInt), 0, 0, 0, 0, 0, 0, 0, 1,
			byte(KindString),
		})), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

}

func TestBadMapGeneric(t *testing.T) {
	var m any

	// bad decoder
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString),
		})), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
		})), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindArray),
		})), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// short
	err = Copy(
		Tokens{
			{Kind: KindMap},
			{Kind: KindString, Value: "foo"},
		}.Iter(), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// short
	err = Copy(
		Tokens{
			{Kind: KindMap},
			{Kind: KindString, Value: "42"},
			{Kind: KindInt8, Value: int8(42)},
		}.Iter(), Unmarshal(

			&m))

	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString), 0,
			byte(KindString),
		})),
		Unmarshal(&m),
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindString), 3, 'F', 'o', 'o',
			byte(KindString),
		})), Unmarshal(&m))
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindInt), 0, 0, 0, 0, 0, 0, 0, 1,
			byte(KindString),
		})), Unmarshal(&m))
	if err == nil {
		t.Fatal(err)
	}

	// bad key
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindMap),
			byte(KindArray),
			byte(KindArrayEnd),
		})), Unmarshal(&m))
	if err == nil {
		t.Fatal(err)
	}

}

func TestUnmarshalDeepRef(t *testing.T) {
	var p ****int
	err := Copy(Tokens{
		{Kind: KindInt, Value: 42},
	}.Iter(), Unmarshal(

		&p))

	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal()
	}
	if ****p != 42 {
		t.Fatal()
	}

	var p2 *int
	err = Copy(Tokens{
		{Kind: KindInt, Value: 42},
	}.Iter(), Unmarshal(

		&p2))

	if err != nil {
		t.Fatal(err)
	}
	if p2 == nil {
		t.Fatal()
	}
	if *p2 != 42 {
		t.Fatal()
	}

	var p3 **int
	err = Copy(Tokens{
		{Kind: KindNil},
	}.Iter(), Unmarshal(

		&p3))

	if err != nil {
		t.Fatal(err)
	}
	if p3 != nil {
		t.Fatal()
	}
}

func TestBadTuple(t *testing.T) {
	var tuple func() int

	// bad token
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindString),
		})), Unmarshal(&tuple))
	if err == nil {
		t.Fatal()
	}

	// short token
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
		})), Unmarshal(&tuple))
	if err == nil {
		t.Fatal()
	}

	// too few items
	err = Copy(Marshal(func() {}), Unmarshal(&tuple))
	if !errors.Is(err, ExpectingValue) {
		t.Fatal()
	}

	// type not match
	err = Copy(
		Marshal(func() (int, int) {
			return 42, 42
		}), Unmarshal(&tuple))
	if !errors.Is(err, BadTupleType) {
		t.Fatal()
	}

	// bad item
	err = Copy(
		Marshal(func() string {
			return "42"
		}), Unmarshal(

			&tuple))

	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindInt), 0, 0, 0, 0, 0, 0, 0, 42,
			byte(KindString),
		})), Unmarshal(

			&tuple))

	if err == nil {
		t.Fatal()
	}

}

func TestBadTupleCall(t *testing.T) {
	var tuple func(int)

	// bad token
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindString),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

	// short token
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

	// too few items
	err = Copy(
		Marshal(func() {}), Unmarshal(

			tuple))

	if !errors.Is(err, ExpectingValue) {
		t.Fatal()
	}

	// too many items
	err = Copy(
		Marshal(func() (int, int) {
			return 42, 42
		}), Unmarshal(
			&tuple))
	if !errors.Is(err, BadTupleType) {
		t.Fatal()
	}

	// too many items
	err = Copy(
		Marshal(func() (int, int) {
			return 42, 42
		}), Unmarshal(
			tuple))
	if !errors.Is(err, BadTupleType) {
		t.Fatal()
	}

	// bad item
	err = Copy(
		Marshal(func() string {
			return "42"
		}), Unmarshal(

			tuple))

	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindInt), 0, 0, 0, 0, 0, 0, 0, 42,
			byte(KindString),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

}

func TestBadTupleCallVariadic(t *testing.T) {
	var tuple func(args ...int)

	// bad token
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindString),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

	// short token
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

	// bad item
	err = Copy(
		Marshal(func() string {
			return "42"
		}), Unmarshal(

			tuple))

	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindInt), 0, 0, 0, 0, 0, 0, 0, 42,
			byte(KindString),
		})), Unmarshal(

			tuple))

	if err == nil {
		t.Fatal()
	}

}

func TestBadGenericTuple(t *testing.T) {
	var tuple any

	// bad token
	err := Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindString),
		})), Unmarshal(

			&tuple))

	if err == nil {
		t.Fatal()
	}

	// short token
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
		})), Unmarshal(

			&tuple))

	if err == nil {
		t.Fatal()
	}

	// bad end
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindBool), 1,
			byte(KindString),
		})), Unmarshal(

			&tuple))

	if err == nil {
		t.Fatal()
	}

	// bad item
	err = Copy(
		Decode(bytes.NewReader([]byte{
			byte(KindTuple),
			byte(KindMapEnd),
		})), Unmarshal(

			&tuple))

	if err == nil {
		t.Fatal()
	}

}

func TestUnmarshalTupleCall(t *testing.T) {
	fn := func(a int, b int) {
		if a != 42 {
			t.Fatal()
		}
		if b != 1 {
			t.Fatal()
		}
	}
	if err := Copy(
		Marshal(
			func() (int, int) {
				return 42, 1
			},
		), Unmarshal(

			fn)); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalTupleCallVariadic(t *testing.T) {
	fn := func(args ...any) {
		if len(args) != 2 {
			t.Fatal()
		}
		if i, ok := args[0].(int); !ok || i != 42 {
			t.Fatal()
		}
		if i, ok := args[1].(int); !ok || i != 1 {
			t.Fatal()
		}
	}
	if err := Copy(
		Marshal(
			func() (int, int) {
				return 42, 1
			},
		), Unmarshal(

			fn)); err != nil {
		t.Fatal(err)
	}
}

func TestBadUnmarshalTarget(t *testing.T) {
	err := Copy(
		Marshal(42), Unmarshal(

			42))

	if !errors.Is(err, BadTargetType) {
		t.Fatal()
	}
}

func TestUnmarshalTupleToErrCaller(t *testing.T) {
	errFoo := fmt.Errorf("foo")
	err := Copy(
		Marshal(func() int {
			return 42
		}),
		Unmarshal(
			func(i int) error {
				if i != 42 {
					t.Fatal()
				}
				return errFoo
			},
		),
	)
	if !is(err, errFoo) {
		t.Fatal()
	}
	if !is(err, UnmarshalError) {
		t.Fatal()
	}
}

func TestUnmarshalTupleToErrVariadicCaller(t *testing.T) {
	errFoo := fmt.Errorf("foo")
	err := Copy(
		Marshal(func() int {
			return 42
		}),
		Unmarshal(
			func(args ...any) error {
				if len(args) != 1 {
					t.Fatal()
				}
				if args[0].(int) != 42 {
					t.Fatal()
				}
				return errFoo
			},
		),
	)
	if !is(err, errFoo) {
		t.Fatal()
	}
	if !is(err, UnmarshalError) {
		t.Fatal()
	}
}

func TestUnmarshalTupleToCallerNoError(t *testing.T) {
	err := Copy(
		Marshal(func() int {
			return 42
		}), Unmarshal(

			func(args ...any) error {
				return nil
			}))

	if err != nil {
		t.Fatal(err)
	}

	err = Copy(
		Marshal(func() int {
			return 42
		}), Unmarshal(

			func(i int) error {
				return nil
			}))

	if err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalToPointer(t *testing.T) {
	var i *int
	err := Copy(
		Marshal(true), Unmarshal(

			&i))

	if err == nil {
		t.Fatal()
	}
	if i != nil {
		t.Fatal()
	}
}

func TestUnmarshalNewImpl(t *testing.T) {
	// unmarshal to any
	{
		var value any
		if err := Copy(
			Marshal(Tuple{
				42, "foo", true,
			}),
			Unmarshal(&value),
		); err != nil {
			t.Fatal(err)
		}
		if fn, ok := value.(func() (int, string, bool)); !ok {
			t.Fatalf("got %T", value)
		} else {
			i, s, b := fn()
			if i != 42 {
				t.Fatal()
			}
			if s != "foo" {
				t.Fatal()
			}
			if !b {
				t.Fatal()
			}
		}
	}

	// unmarshal to func call
	{
		if err := Copy(
			Marshal(Tuple{
				42, "foo", true,
			}),
			Unmarshal(func(i int, s string, b bool) {
				if i != 42 {
					t.Fatal()
				}
				if s != "foo" {
					t.Fatal()
				}
				if !b {
					t.Fatal()
				}
			}),
		); err != nil {
			t.Fatal()
		}
	}

	// unmarshal to tuple func
	{
		var fn func() (int, string, bool)
		if err := Copy(
			Marshal(Tuple{
				42, "foo", true,
			}),
			Unmarshal(&fn),
		); err != nil {
			t.Fatal()
		}
		i, s, b := fn()
		if i != 42 {
			t.Fatal()
		}
		if s != "foo" {
			t.Fatal()
		}
		if !b {
			t.Fatal()
		}
	}

	// unmarshal to ellipses
	{
		if err := Copy(
			Marshal(Tuple{
				42, "foo", true,
			}),
			Unmarshal(func(tuple ...any) {
				if len(tuple) != 3 {
					t.Fatal()
				}
				if i, ok := tuple[0].(int); !ok || i != 42 {
					t.Fatal()
				}
				if s, ok := tuple[1].(string); !ok || s != "foo" {
					t.Fatal()
				}
				if b, ok := tuple[2].(bool); !ok || !b {
					t.Fatal()
				}
			}),
		); err != nil {
			t.Fatal()
		}
	}

}

func TestUnmarshalStructUnknownField(t *testing.T) {
	var s struct {
		A int
	}
	err := Copy(
		Marshal(struct {
			B int
		}{}),
		UnmarshalValue(
			Ctx{}.Strict(),
			reflect.ValueOf(&s),
			nil,
		),
	)
	if !is(err, UnknownFieldName) {
		t.Fatal()
	}
}

func TestUnmarshalTupleFunc(t *testing.T) {
	var fn func() (int, int, int)
	if err := Copy(
		Marshal(Tuple{1, 2, 3}),
		Unmarshal(&fn),
	); err != nil {
		t.Fatal(err)
	}
	a, b, c := fn()
	if a != 1 {
		t.Fatal()
	}
	if b != 2 {
		t.Fatal()
	}
	if c != 3 {
		t.Fatal()
	}
}

func TestUnmarshalEmbedded(t *testing.T) {
	type Foo struct {
		time.Time
	}
	var f Foo
	now := time.Now()
	if err := Copy(
		Marshal(Foo{
			Time: now,
		}),
		Unmarshal(&f),
	); err != nil {
		t.Fatal()
	}
	if !f.Time.Equal(now) {
		t.Fatalf("got %v, expected %v", f.Time, now)
	}
}

func TestUnmarshalToEmbedded(t *testing.T) {
	type Foo struct {
		Bar int
	}
	var data struct {
		Foo
	}
	if err := Copy(
		Marshal(struct {
			Bar int
		}{
			Bar: 42,
		}),
		Unmarshal(&data),
	); err != nil {
		t.Fatal()
	}
	if data.Bar != 42 {
		t.Fatal()
	}
}

func TestUnmarshalPath(t *testing.T) {

	// slice
	var ints []int
	if err := Copy(
		Marshal([]int{1, 2, 3}),
		TapUnmarshal(DefaultCtx, &ints, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindArray {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Kind == KindInt {
				if token.Value == 1 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 0 {
						t.Fatal()
					}
				} else if token.Value == 2 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 1 {
						t.Fatal()
					}
				} else if token.Value == 3 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 2 {
						t.Fatal()
					}
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// generic slice
	var value any
	if err := Copy(
		Marshal([]int{1, 2, 3}),
		TapUnmarshal(DefaultCtx, &value, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindArray {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Kind == KindInt {
				if token.Value == 1 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 0 {
						t.Fatal()
					}
				} else if token.Value == 2 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 1 {
						t.Fatal()
					}
				} else if token.Value == 3 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 2 {
						t.Fatal()
					}
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// array
	var array []int
	if err := Copy(
		Marshal([]int{1, 2, 3}),
		TapUnmarshal(DefaultCtx, &array, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindArray {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Kind == KindInt {
				if token.Value == 1 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 0 {
						t.Fatal()
					}
				} else if token.Value == 2 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 1 {
						t.Fatal()
					}
				} else if token.Value == 3 {
					if len(ctx.Path) != 1 {
						t.Fatal()
					}
					if ctx.Path[0] != 2 {
						t.Fatal()
					}
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// struct
	type foo struct {
		Foo int
		Bar string
	}
	var f foo
	if err := Copy(
		Marshal(foo{
			Foo: 42,
			Bar: "bar",
		}),
		TapUnmarshal(DefaultCtx, &f, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindObject {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if token.Value == "bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// skip struct field
	if err := Copy(
		Marshal(struct {
			Foo int
			Bar string
			Baz int
		}{
			Foo: 42,
			Bar: "bar",
			Baz: 1,
		}),
		TapUnmarshal(DefaultCtx, &f, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindObject {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if token.Value == "bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			} else if token.Value == 1 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Baz" {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// new struct
	if err := Copy(
		Marshal(foo{
			Foo: 42,
			Bar: "bar",
		}),
		TapUnmarshal(DefaultCtx, &value, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindObject {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if token.Value == "bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// map
	var m map[any]any
	if err := Copy(
		Marshal(map[any]any{
			"Foo": 42,
			"Bar": "bar",
		}),
		TapUnmarshal(DefaultCtx, &m, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindMap {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if token.Value == "bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// generic map
	if err := Copy(
		Marshal(map[any]any{
			"Foo": 42,
			"Bar": "bar",
		}),
		TapUnmarshal(DefaultCtx, &value, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindMap {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if token.Value == "bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// tuple
	var tuple func() (int, string)
	if err := Copy(
		Marshal(func() (int, string) {
			return 42, "foo"
		}),
		TapUnmarshal(DefaultCtx, &tuple, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindTuple {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if token.Value == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// generic tuple
	if err := Copy(
		Marshal(func() (int, string) {
			return 42, "foo"
		}),
		TapUnmarshal(DefaultCtx, &value, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindTuple {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if token.Value == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// Tuple
	var tuple2 Tuple
	if err := Copy(
		Marshal(func() (int, string) {
			return 42, "foo"
		}),
		TapUnmarshal(DefaultCtx, &tuple2, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindTuple {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if token.Value == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

	// TypedTuple
	tuple3 := TypedTuple{
		Types: TupleTypes(func(int, string) {}),
	}
	if err := Copy(
		Marshal(func() (int, string) {
			return 42, "foo"
		}),
		TapUnmarshal(DefaultCtx, &tuple3, func(ctx Ctx, token Token, target reflect.Value) {
			if token.Kind == KindTuple {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			}
			if token.Value == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if token.Value == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
	); err != nil {
		t.Fatal(err)
	}

}

func TestTapUnmarshal(t *testing.T) {
	ok := false
	err := Copy(
		Tokens{}.Iter(),
		TapUnmarshal(DefaultCtx, nil, func(ctx Ctx, token Token, target reflect.Value) {
			ok = true
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}
}

func TestUnmarshalByteArrayKeyMap(t *testing.T) {
	var v any
	if err := Copy(
		Marshal(map[[2]byte]bool{
			{1, 2}: true,
		}),
		Unmarshal(&v),
	); err != nil {
		t.Fatal(err)
	}
	m, ok := v.(map[any]any)
	if !ok {
		t.Fatal()
	}
	value, ok := m[[2]byte{1, 2}]
	if !ok {
		t.Fatal()
	}
	if value != true {
		t.Fatal()
	}
}

func TestUnmarshalWithTypeName(t *testing.T) {
	var v any
	if err := Copy(
		Marshal(definedInt(42)),
		Unmarshal(&v),
	); err != nil {
		t.Fatal(err)
	}
	if v != definedInt(42) {
		t.Fatal()
	}
	if _, ok := v.(definedInt); !ok {
		t.Fatal()
	}
}
