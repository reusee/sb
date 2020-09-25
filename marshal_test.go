package sb

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"
)

type MarshalTestCase struct {
	value    any
	expected []Token
}

type foo int

var marshalTestCases = []MarshalTestCase{
	0: {
		int(42),
		[]Token{
			{Kind: KindInt, Value: int(42)},
		},
	},

	1: {
		func() *int32 {
			i := int32(42)
			return &i
		}(),
		[]Token{
			{Kind: KindInt32, Value: int32(42)},
		},
	},

	2: {
		true,
		[]Token{
			{Kind: KindBool, Value: true},
		},
	},

	3: {
		uint32(42),
		[]Token{
			{Kind: KindUint32, Value: uint32(42)},
		},
	},

	4: {
		float32(42),
		[]Token{
			{Kind: KindFloat32, Value: float32(42)},
		},
	},

	5: {
		[]int{42, 4, 2},
		[]Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindInt, Value: int(4)},
			{Kind: KindInt, Value: int(2)},
			{Kind: KindArrayEnd},
		},
	},

	6: {
		[][]int{
			{42, 4, 2},
			{2, 4, 42},
		},
		[]Token{
			{Kind: KindArray},
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindInt, Value: int(4)},
			{Kind: KindInt, Value: int(2)},
			{Kind: KindArrayEnd},
			{Kind: KindArray},
			{Kind: KindInt, Value: int(2)},
			{Kind: KindInt, Value: int(4)},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindArrayEnd},
			{Kind: KindArrayEnd},
		},
	},

	7: {
		"foo",
		[]Token{
			{Kind: KindString, Value: "foo"},
		},
	},

	8: {
		struct {
			Foo int
			Bar float32
			Baz string
			Boo bool
		}{
			42,
			42,
			"42",
			false,
		},
		[]Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Bar"},
			{Kind: KindFloat32, Value: float32(42)},
			{Kind: KindString, Value: "Baz"},
			{Kind: KindString, Value: "42"},
			{Kind: KindString, Value: "Boo"},
			{Kind: KindBool, Value: false},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindObjectEnd},
		},
	},

	9: {
		[]interface{}{
			func() **int {
				i := 42
				j := &i
				return &j
			}(),
			func() *int {
				return nil
			}(),
			(interface{})(nil),
		},
		[]Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindNil},
			{Kind: KindNil},
			{Kind: KindArrayEnd},
		},
	},

	10: {
		func() *int {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	11: {
		func() *bool {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	12: {
		func() *uint8 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	13: {
		func() *float32 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	14: {
		func() *[]int32 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	15: {
		func() *[]int32 {
			array := []int32{42}
			return &array
		}(),
		[]Token{
			{Kind: KindArray},
			{Kind: KindInt32, Value: int32(42)},
			{Kind: KindArrayEnd},
		},
	},

	16: {
		func() *string {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	17: {
		func() *struct{} {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	18: func() MarshalTestCase {
		str := strings.Repeat("foo", 1024)
		return MarshalTestCase{
			str,
			[]Token{
				{Kind: KindString, Value: str},
			},
		}
	}(),

	19: {
		foo(42),
		[]Token{
			{Kind: KindInt, Value: int(42)},
		},
	},

	20: {
		[]foo{
			42,
		},
		[]Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindArrayEnd},
		},
	},

	21: {
		struct {
			Foo foo
		}{
			42,
		},
		[]Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindObjectEnd},
		},
	},

	22: {
		struct {
			Foo []foo
		}{
			[]foo{
				42,
			},
		},
		[]Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindArrayEnd},
			{Kind: KindObjectEnd},
		},
	},

	23: {
		map[int]int{},
		[]Token{
			{Kind: KindMap},
			{Kind: KindMapEnd},
		},
	},

	24: {
		map[int]int{
			1: 1,
		},
		[]Token{
			{Kind: KindMap},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 1},
			{Kind: KindMapEnd},
		},
	},

	25: {
		map[int]int{
			42: 42,
			1:  1,
		},
		[]Token{
			{Kind: KindMap},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 42},
			{Kind: KindInt, Value: 42},
			{Kind: KindMapEnd},
		},
	},

	26: {
		map[any]any{
			42:    42,
			"foo": "bar",
		},
		[]Token{
			{Kind: KindMap},
			{Kind: KindString, Value: "foo"},
			{Kind: KindString, Value: "bar"},
			{Kind: KindInt, Value: 42},
			{Kind: KindInt, Value: 42},
			{Kind: KindMapEnd},
		},
	},

	27: {
		func() (int, string) {
			return 42, "42"
		},
		[]Token{
			{Kind: KindTuple},
			{Kind: KindInt, Value: 42},
			{Kind: KindString, Value: "42"},
			{Kind: KindTupleEnd},
		},
	},

	28: {
		marshalStringAsInt("foo"),
		[]Token{
			{Kind: KindInt, Value: int(3)},
		},
	},
}

func TestMarshaler(t *testing.T) {

	for i, c := range marshalTestCases {

		// marshal
		tokens, err := TokensFromStream(Marshal(c.value))
		if err != nil {
			t.Fatal(err)
		}
		if len(tokens) != len(c.expected) {
			t.Fatalf("%d fail %+v", i, c)
		}
		for idx, token := range tokens {
			if token != c.expected[idx] {
				pt("expected %T\n", c.expected[idx].Value)
				pt("token %T\n", token.Value)
				t.Fatalf(
					"case %d token %d\nexpected %#v\ngot %#v\nfail %+v",
					i, idx,
					c.expected[idx], token, c,
				)
			}
		}

		// encode
		buf := new(bytes.Buffer)
		if err := Copy(
			Marshal(c.value),
			Encode(buf),
		); err != nil {
			t.Fatal(err)
		}
		decoder := Decode(buf)
		if MustCompare(decoder, Tokens(c.expected).Iter()) != 0 {
			t.Fatalf("%d fail %+v", i, c)
		}

		// compare
		tokens, err = TokensFromStream(Marshal(c.value))
		if err != nil {
			t.Fatal(err)
		}
		var obj any
		if err := Copy(tokens.Iter(), Unmarshal(&obj)); err != nil {
			t.Fatal(err)
		}
		if MustCompare(Marshal(obj), Marshal(c.value)) != 0 {
			t.Fatalf("not equal, got %#v, expected %#v", obj, c.value)
		}

	}

}

type Custom struct {
	Foo int
}

var _ SBMarshaler = Custom{}

var _ SBUnmarshaler = new(Custom)

func (c Custom) MarshalSB(ctx Ctx, cont Proc) Proc {
	return ctx.Marshal(ctx, reflect.ValueOf(c.Foo), cont)
}

func (c *Custom) UnmarshalSB(ctx Ctx, cont Sink) Sink {
	return func(p *Token) (Sink, error) {
		if p == nil {
			return cont, nil
		}
		token := *p
		if token.Kind != KindInt {
			return cont, nil
		}
		c.Foo = token.Value.(int)
		return cont, nil
	}
}

func TestCustomType(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(Custom{42}),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}
	var c Custom
	if err := Copy(Decode(buf), Unmarshal(&c)); err != nil {
		t.Fatal(err)
	}
	if c.Foo != 42 {
		t.Fatal()
	}
}

type badBinaryMarshaler struct{}

var _ encoding.BinaryMarshaler = badBinaryMarshaler{}

func (_ badBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("bad")
}

func TestBadBinaryMarshaler(t *testing.T) {
	v := new(badBinaryMarshaler)
	m := Marshal(v)
	_, err := TokensFromStream(m)
	if err == nil {
		t.Fatal()
	}
}

type badTextMarshaler struct{}

var _ encoding.TextMarshaler = badTextMarshaler{}

func (_ badTextMarshaler) MarshalText() ([]byte, error) {
	return nil, fmt.Errorf("bad")
}

func TestBadTextMarshaler(t *testing.T) {
	v := new(badTextMarshaler)
	m := Marshal(v)
	_, err := TokensFromStream(m)
	if err == nil {
		t.Fatal()
	}
}

type timeTextMarshaler struct {
	t time.Time
}

var _ encoding.TextMarshaler = timeTextMarshaler{}

func (t timeTextMarshaler) MarshalText() ([]byte, error) {
	return t.t.MarshalText()
}

var _ encoding.TextUnmarshaler = new(timeTextMarshaler)

func (t *timeTextMarshaler) UnmarshalText(data []byte) error {
	return t.t.UnmarshalText(data)
}

func TestTimeMarshalText(t *testing.T) {
	now := timeTextMarshaler{time.Now()}
	m := Marshal(now)
	tokens, err := TokensFromStream(m)
	if err != nil {
		t.Fatal(err)
	}
	var tt timeTextMarshaler
	if err := Copy(tokens.Iter(), Unmarshal(&tt)); err != nil {
		t.Fatal(err)
	}
	if time.Since(tt.t) > time.Second {
		t.Fatal()
	}
}

func TestMarshalIgnoreUnsupportedType(t *testing.T) {
	tokens, err := TokensFromStream(
		Marshal(
			make(chan int),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 0 {
		t.Fatal()
	}
}

func TestBadMapKey(t *testing.T) {
	_, err := TokensFromStream(
		Marshal(
			map[any]any{
				badBinaryMarshaler{}: true,
			},
		),
	)
	if err == nil {
		t.Fatal()
	}

	_, err = TokensFromStream(
		Marshal(
			map[any]any{
				unsafe.Pointer(nil): true,
			},
		),
	)
	if err == nil {
		t.Fatal()
	}
}

type marshalStringAsInt string

var _ SBMarshaler = marshalStringAsInt("")

func (m marshalStringAsInt) MarshalSB(ctx Ctx, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		return nil, MarshalAny(ctx, len(m), cont), nil
	}
}

func TestValueMarshalFunc(t *testing.T) {
	fn := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		return MarshalValue(ctx, value, cont)
	}
	for _, c := range marshalTestCases {
		proc := fn(Ctx{
			Marshal: fn,
		}, reflect.ValueOf(c.value), nil)
		stream := &proc
		if MustCompare(stream, Tokens(c.expected).Iter()) != 0 {
			t.Fatal("not equal")
		}
	}
}

func TestValueMarshalFunc2(t *testing.T) {
	n := 0
	fn := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		if value.Kind() == reflect.Struct {
			n++
		}
		return MarshalValue(ctx, value, cont)
	}
	proc := fn(Ctx{Marshal: fn}, reflect.ValueOf(
		struct {
			Foo struct {
				Bar struct {
				}
			}
		}{},
	), nil)
	stream := &proc
	_, err := TokensFromStream(stream)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatal()
	}
}

func TestMarshalStructFieldOrder(t *testing.T) {
	expectedKinds := []reflect.Kind{
		reflect.Struct,
		reflect.String,
		reflect.Int,
		reflect.Ptr, // *Token
	}
	fn := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		if value.Kind() != expectedKinds[0] {
			t.Fatal()
		}
		expectedKinds = expectedKinds[1:]
		return MarshalValue(ctx, value, cont)
	}
	proc := fn(Ctx{Marshal: fn}, reflect.ValueOf(struct {
		I int
	}{
		I: 42,
	}), nil)
	if err := Copy(
		&proc,
		Discard,
	); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalMapOrder(t *testing.T) {
	expecteds := []int64{
		1, 2, 3, 4, 5, 6,
	}
	fn := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		if value.Kind() == reflect.Int64 {
			if value.Int() != expecteds[0] {
				t.Fatalf("expected %d, got %d", expecteds[0], value.Int())
			}
			expecteds = expecteds[1:]
		}
		return MarshalValue(ctx, value, cont)
	}
	proc := fn(Ctx{Marshal: fn}, reflect.ValueOf(map[int64]int64{
		1: 2,
		3: 4,
		5: 6,
	}), nil)
	if err := Copy(
		&proc,
		Discard,
	); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalUnexportedField(t *testing.T) {
	if err := Copy(
		Marshal(struct {
			f int
		}{
			f: 42,
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalNonEmpty(t *testing.T) {
	proc1 := MarshalValue(Ctx{
		SkipEmptyStructFields: true,
	}, reflect.ValueOf(struct {
		Foo  int
		Bar  string
		Ints []int
	}{
		Foo:  42,
		Ints: []int{},
	}), nil)
	proc2 := MarshalValue(Ctx{
		Marshal:               MarshalValue,
		SkipEmptyStructFields: true,
	}, reflect.ValueOf(struct {
		Foo int
		Baz bool
	}{
		Foo: 42,
	}), nil)
	if MustCompare(
		&proc1,
		&proc2,
	) != 0 {
		t.Fatal()
	}
}

func TestStructAsMapKey(t *testing.T) {
	m := map[interface{}]interface{}{
		struct {
			X struct {
				P int16
				X struct{}
			}
		}{X: struct {
			P int16
			X struct{}
		}{P: -16693, X: struct{}{}}}: nil,
		struct {
			X struct {
				X struct{}
				P int16
			}
		}{X: struct {
			X struct{}
			P int16
		}{X: struct{}{}, P: -16693}}: nil,
	}

	var o map[interface{}]interface{}
	if err := Copy(
		Marshal(m),
		Unmarshal(&o),
	); err != nil {
		panic(err)
	}
	if len(o) != 2 {
		t.Fatal()
	}

}

func TestMarshalTime(t *testing.T) {
	buf := new(bytes.Buffer)
	now := time.Now()
	if err := Copy(
		Marshal(now),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}
	var tt time.Time
	if err := Copy(Decode(buf), Unmarshal(&tt)); err != nil {
		t.Fatal(err)
	}
	if !tt.Equal(now) {
		t.Fatal()
	}
}

func TestCyclicPointer(t *testing.T) {
	type P *P
	var p P
	p = &p
	err := Copy(
		Marshal(p),
		Discard,
	)
	if !is(err, CyclicPointer) {
		t.Fatal()
	}
}
