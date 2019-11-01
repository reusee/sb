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
}

func TestUnmarshal(t *testing.T) {
	for _, c := range unmarshalTestCases {
		stream := NewMarshaler(c.value)
		ptr := reflect.New(reflect.TypeOf(c.target))
		err := Unmarshal(stream, ptr.Interface())
		if !errors.Is(err, c.err) {
			t.Fatal()
		}
		if err == nil {
			if !reflect.DeepEqual(c.value, ptr.Elem().Interface()) {
				t.Fatal()
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
	}

	for _, c := range cases {
		var v any
		err := Unmarshal(c.Iter(), &v)
		if err == nil {
			t.Fatal()
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
