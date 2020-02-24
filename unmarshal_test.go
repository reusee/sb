package sb

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
)

type UnmarshalTestCase struct {
	value  any
	target any
	err    error
}

type testBytes []byte

var unmarshalTestCases = []UnmarshalTestCase{
	{true, bool(false), nil},
	{int8(42), int8(0), nil},
	{int16(42), int16(0), nil},
	{int32(42), int32(0), nil},
	{int64(42), int64(0), nil},
	{uint(42), uint(0), nil},
	{uint8(42), uint8(0), nil},
	{uint16(42), uint16(0), nil},
	{uint32(42), uint32(0), nil},
	{uint64(42), uint64(0), nil},
	{float32(42), float32(0), nil},
	{float64(42), float64(0), nil},
	{string("42"), string(""), nil},
	{map[int]int{1: 1}, map[int]int(nil), nil},
	{[]byte("foo"), []byte("foo"), nil},
	{testBytes{42}, testBytes{42}, nil},
	{[3]int{1}, [3]int{1}, nil},
	{
		func() (int, string) {
			return 42, "42"
		},
		(func() (int, string))(nil),
		nil,
	},

	{true, int(0), ExpectingBool},
	{42, true, ExpectingInt},
	{int8(42), true, ExpectingInt8},
	{int16(42), true, ExpectingInt16},
	{int32(42), true, ExpectingInt32},
	{int64(42), true, ExpectingInt64},
	{uint(42), true, ExpectingUint},
	{uint8(42), true, ExpectingUint8},
	{uint16(42), true, ExpectingUint16},
	{uint32(42), true, ExpectingUint32},
	{uint64(42), true, ExpectingUint64},
	{float32(42), true, ExpectingFloat32},
	{float64(42), true, ExpectingFloat64},
	{math.NaN(), true, ExpectingFloat},
	{"42", true, ExpectingString},
	{[]int{42}, true, ExpectingSequence},
	{map[int]int{}, true, ExpectingMap},
	{42, (****string)(nil), ExpectingInt},
	{[]byte("foo"), true, ExpectingBytes},
	{
		func() int { return 42 },
		true,
		ExpectingTuple,
	},
	{
		func() int { return 42 },
		(func(int) int)(nil),
		BadTupleType,
	},
}

func TestUnmarshal(t *testing.T) {
	for _, c := range unmarshalTestCases {
		stream := NewMarshaler(c.value)
		ptr := reflect.New(reflect.TypeOf(c.target))
		err := Unmarshal(stream, ptr.Interface())
		if !errors.Is(err, c.err) {
			pt("%v\n", err)
			t.Fatal()
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
					t.Fatal()
				}
			}
		}
	}
}

func TestUnmarshalNaN(t *testing.T) {
	stream := NewMarshaler(math.NaN())
	var f float64
	if err := Unmarshal(stream, &f); err != nil {
		t.Fatal(err)
	}
	if !math.IsNaN(f) {
		t.Fatal()
	}
}

func TestUnmarshalArray(t *testing.T) {
	type foo [3]byte
	type S struct {
		Foos []foo
	}
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(S{
		Foos: []foo{
			foo{1},
			foo{2},
		},
	})); err != nil {
		t.Fatal(err)
	}
	var s S
	if err := Unmarshal(NewDecoder(buf), &s); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalNamedUint(t *testing.T) {
	type Foo uint32
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(Foo(42))); err != nil {
		t.Fatal(err)
	}
	var foo Foo
	if err := Unmarshal(NewDecoder(buf), &foo); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalStructWithPrivateField(t *testing.T) {
	type Foo struct {
		Bar int
		Foo int
	}
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(Foo{42, 42})); err != nil {
		t.Fatal(err)
	}
	type Bar struct {
		bar int
		Foo int
	}
	var bar Bar
	if err := Unmarshal(NewDecoder(buf), &bar); err != nil {
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
			{KindString, "Foo"},
		},
		{
			{Kind: KindObject},
			{KindString, "Foo"},
			{KindString, "Bar"},
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
		err := Unmarshal(c.Iter(), &v)
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
	if err := Unmarshal(NewDecoder(bytes.NewReader([]byte{
		KindString,
	})), &b); err == nil {
		t.Fatal()
	}

	// no token
	if err := Unmarshal(NewDecoder(bytes.NewReader(nil)), &b); err == nil {
		t.Fatal()
	}

	// bad token
	if err := Unmarshal(Tokens{
		{KindInt, 42},
	}.Iter(), &b); err == nil {
		t.Fatal()
	}

	// bad unmarshaler
	if err := Unmarshal(Tokens{
		{KindString, "foo"},
	}.Iter(), &b); err == nil {
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
	if err := Unmarshal(NewDecoder(bytes.NewReader([]byte{
		KindString,
	})), &b); err == nil {
		t.Fatal()
	}

	// no token
	if err := Unmarshal(NewDecoder(bytes.NewReader(nil)), &b); err == nil {
		t.Fatal()
	}

	// bad token
	if err := Unmarshal(Tokens{
		{KindInt, 42},
	}.Iter(), &b); err == nil {
		t.Fatal()
	}

	// bad unmarshaler
	if err := Unmarshal(Tokens{
		{KindString, "foo"},
	}.Iter(), &b); err == nil {
		t.Fatal()
	}

}

func TestUnmarshalToNilPtr(t *testing.T) {
	if err := Unmarshal(Tokens{
		{KindInt, 42},
	}.Iter(), (*int)(nil)); err != nil {
		t.Fatal(err)
	}
}

func TestBadArray(t *testing.T) {
	var v [2]int

	// bad decoder
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindArray,
			KindString,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindArray,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// too many element
	err = Unmarshal(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{KindInt, 42},
			{KindInt, 42},
			{KindInt, 42},
		}.Iter(),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{KindInt8, int8(42)},
		}.Iter(),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}
}

func TestBadSlice(t *testing.T) {
	var v []int

	// bad decoder
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindArray,
			KindString,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindArray,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindArray},
			{KindInt, 42},
			{KindInt8, int8(42)},
		}.Iter(),
		&v,
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
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindObject,
			KindString,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindObject,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindObject},
			{KindInt, 42},
		}.Iter(),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindObject},
			{KindString, "42"},
			{KindInt8, int8(42)},
		}.Iter(),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindObject,
			KindString, 0,
			KindString,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindObject,
			KindString, 3, 'F', 'o', 'o',
			KindString,
		})),
		&v,
	)
	if err == nil {
		t.Fatal(err)
	}

}

func TestBadMap(t *testing.T) {
	var m map[int]int

	// bad decoder
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindMap},
			{KindString, "foo"},
		}.Iter(),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad type
	err = Unmarshal(
		Tokens{
			{Kind: KindMap},
			{KindString, "42"},
			{KindInt8, int8(42)},
		}.Iter(),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString, 0,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString, 3, 'F', 'o', 'o',
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindInt, 0, 0, 0, 0, 0, 0, 0, 1,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

}

func TestBadMapGeneric(t *testing.T) {
	var m any

	// bad decoder
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short decoder
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindArray,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short
	err = Unmarshal(
		Tokens{
			{Kind: KindMap},
			{KindString, "foo"},
		}.Iter(),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// short
	err = Unmarshal(
		Tokens{
			{Kind: KindMap},
			{KindString, "42"},
			{KindInt8, int8(42)},
		}.Iter(),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad skip
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString, 0,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindString, 3, 'F', 'o', 'o',
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad value
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindInt, 0, 0, 0, 0, 0, 0, 0, 1,
			KindString,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

	// bad key
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindMap,
			KindArray,
			KindArrayEnd,
		})),
		&m,
	)
	if err == nil {
		t.Fatal(err)
	}

}

func TestUnmarshalDeepRef(t *testing.T) {
	var p ****int
	err := Unmarshal(Tokens{
		{KindInt, 42},
	}.Iter(), &p)
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
	err = Unmarshal(Tokens{
		{KindInt, 42},
	}.Iter(), &p2)
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
	err = Unmarshal(Tokens{
		{Kind: KindNil},
	}.Iter(), &p3)
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
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindString,
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// short token
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// too few items
	err = Unmarshal(
		NewMarshaler(func() {}),
		&tuple,
	)
	if !errors.Is(err, ExpectingValue) {
		t.Fatal()
	}

	// too many items
	err = Unmarshal(
		NewMarshaler(func() (int, int) {
			return 42, 42
		}),
		&tuple,
	)
	if !errors.Is(err, TooManyElement) {
		t.Fatal()
	}

	// bad item
	err = Unmarshal(
		NewMarshaler(func() string {
			return "42"
		}),
		&tuple,
	)
	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindInt, 0, 0, 0, 0, 0, 0, 0, 42,
			KindString, // incomplete string
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

}

func TestBadTupleCall(t *testing.T) {
	var tuple func(int)

	// bad token
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindString,
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// short token
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// too few items
	err = Unmarshal(
		NewMarshaler(func() {}),
		tuple,
	)
	if !errors.Is(err, ExpectingValue) {
		t.Fatal()
	}

	// too many items
	err = Unmarshal(
		NewMarshaler(func() (int, int) {
			return 42, 42
		}),
		tuple,
	)
	if !errors.Is(err, TooManyElement) {
		t.Fatal()
	}

	// bad item
	err = Unmarshal(
		NewMarshaler(func() string {
			return "42"
		}),
		tuple,
	)
	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindInt, 0, 0, 0, 0, 0, 0, 0, 42,
			KindString, // incomplete string
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

}

func TestBadTupleCallVariadic(t *testing.T) {
	var tuple func(args ...int)

	// bad token
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindString,
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// short token
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// bad item
	err = Unmarshal(
		NewMarshaler(func() string {
			return "42"
		}),
		tuple,
	)
	if !errors.Is(err, ExpectingString) {
		t.Fatal()
	}

	// bad end
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindInt, 0, 0, 0, 0, 0, 0, 0, 42,
			KindString, // incomplete string
		})),
		tuple,
	)
	if err == nil {
		t.Fatal()
	}

}

func TestBadGenericTuple(t *testing.T) {
	var tuple any

	// bad token
	err := Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindString,
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// short token
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// bad end
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindBool, 1,
			KindString, // incomplete string
		})),
		&tuple,
	)
	if err == nil {
		t.Fatal()
	}

	// bad item
	err = Unmarshal(
		NewDecoder(bytes.NewReader([]byte{
			KindTuple,
			KindMapEnd,
		})),
		&tuple,
	)
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
	if err := Unmarshal(
		NewMarshaler(
			func() (int, int) {
				return 42, 1
			},
		),
		fn,
	); err != nil {
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
	if err := Unmarshal(
		NewMarshaler(
			func() (int, int) {
				return 42, 1
			},
		),
		fn,
	); err != nil {
		t.Fatal(err)
	}
}

func TestBadUnmarshalTarget(t *testing.T) {
	err := Unmarshal(
		NewMarshaler(42),
		42,
	)
	if !errors.Is(err, BadTargetType) {
		t.Fatal()
	}
}

func TestUnmarshalTupleToErrCaller(t *testing.T) {
	err := Unmarshal(
		NewMarshaler(func() int {
			return 42
		}),
		func(i int) error {
			if i != 42 {
				t.Fatal()
			}
			return fmt.Errorf("foo")
		},
	)
	if err.Error() != "foo" {
		t.Fatal()
	}
}

func TestUnmarshalTupleToErrVariadicCaller(t *testing.T) {
	err := Unmarshal(
		NewMarshaler(func() int {
			return 42
		}),
		func(args ...any) error {
			if len(args) != 1 {
				t.Fatal()
			}
			if args[0].(int) != 42 {
				t.Fatal()
			}
			return fmt.Errorf("foo")
		},
	)
	if err.Error() != "foo" {
		t.Fatal()
	}
}

func TestUnmarshalTupleToCallerNoError(t *testing.T) {
	err := Unmarshal(
		NewMarshaler(func() int {
			return 42
		}),
		func(args ...any) error {
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = Unmarshal(
		NewMarshaler(func() int {
			return 42
		}),
		func(i int) error {
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalToPointer(t *testing.T) {
	var i *int
	err := Unmarshal(
		NewMarshaler(true),
		&i,
	)
	if err == nil {
		t.Fatal()
	}
	if i != nil {
		t.Fatal()
	}
}
