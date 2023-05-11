package sb

import (
	"bytes"
	"encoding"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"
)

type MarshalTestCase struct {
	value    any
	expected []Token
	ctx      Ctx
}

type foo int

type testBool bool

type testInt8 int8

type testInt16 int16

type testInt32 int32

type testInt64 int64

type testUint uint

type testUint8 uint8

type testUint16 uint16

type testUint64 uint64

type testFloat32 float32

type testFloat64 float64

type testString string

type definedInt int

func init() {
	Register(reflect.TypeOf((*definedInt)(nil)).Elem())
}

var marshalTestCases = []MarshalTestCase{
	0: {
		value: int(42),
		expected: []Token{
			{Kind: KindInt, Value: int(42)},
		},
	},

	1: {
		value: func() *int32 {
			i := int32(42)
			return &i
		}(),
		expected: []Token{
			{Kind: KindInt32, Value: int32(42)},
		},
	},

	2: {
		value: true,
		expected: []Token{
			{Kind: KindBool, Value: true},
		},
	},

	3: {
		value: uint32(42),
		expected: []Token{
			{Kind: KindUint32, Value: uint32(42)},
		},
	},

	4: {
		value: float32(42),
		expected: []Token{
			{Kind: KindFloat32, Value: float32(42)},
		},
	},

	5: {
		value: []int{42, 4, 2},
		expected: []Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindInt, Value: int(4)},
			{Kind: KindInt, Value: int(2)},
			{Kind: KindArrayEnd},
		},
	},

	6: {
		value: [][]int{
			{42, 4, 2},
			{2, 4, 42},
		},
		expected: []Token{
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
		value: "foo",
		expected: []Token{
			{Kind: KindString, Value: "foo"},
		},
	},

	8: {
		value: struct {
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
		expected: []Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindString, Value: "Bar"},
			{Kind: KindFloat32, Value: float32(42)},
			{Kind: KindString, Value: "Baz"},
			{Kind: KindString, Value: "42"},
			{Kind: KindString, Value: "Boo"},
			{Kind: KindBool, Value: false},
			{Kind: KindObjectEnd},
		},
	},

	9: {
		value: []interface{}{
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
		expected: []Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindNil},
			{Kind: KindNil},
			{Kind: KindArrayEnd},
		},
	},

	10: {
		value: func() *int {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	11: {
		value: func() *bool {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	12: {
		value: func() *uint8 {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	13: {
		value: func() *float32 {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	14: {
		value: func() *[]int32 {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	15: {
		value: func() *[]int32 {
			array := []int32{42}
			return &array
		}(),
		expected: []Token{
			{Kind: KindArray},
			{Kind: KindInt32, Value: int32(42)},
			{Kind: KindArrayEnd},
		},
	},

	16: {
		value: func() *string {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	17: {
		value: func() *struct{} {
			return nil
		}(),
		expected: []Token{
			{Kind: KindNil},
		},
	},

	18: func() MarshalTestCase {
		str := strings.Repeat("foo", 1024)
		return MarshalTestCase{
			value: str,
			expected: []Token{
				{Kind: KindString, Value: str},
			},
		}
	}(),

	19: {
		value: foo(42),
		expected: []Token{
			{Kind: KindInt, Value: int(42)},
		},
	},

	20: {
		value: []foo{
			42,
		},
		expected: []Token{
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindArrayEnd},
		},
	},

	21: {
		value: struct {
			Foo foo
		}{
			42,
		},
		expected: []Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindObjectEnd},
		},
	},

	22: {
		value: struct {
			Foo []foo
		}{
			[]foo{
				42,
			},
		},
		expected: []Token{
			{Kind: KindObject},
			{Kind: KindString, Value: "Foo"},
			{Kind: KindArray},
			{Kind: KindInt, Value: int(42)},
			{Kind: KindArrayEnd},
			{Kind: KindObjectEnd},
		},
	},

	23: {
		value: map[int]int{},
		expected: []Token{
			{Kind: KindMap},
			{Kind: KindMapEnd},
		},
	},

	24: {
		value: map[int]int{
			1: 1,
		},
		expected: []Token{
			{Kind: KindMap},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 1},
			{Kind: KindMapEnd},
		},
	},

	25: {
		value: map[int]int{
			42: 42,
			1:  1,
		},
		expected: []Token{
			{Kind: KindMap},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 1},
			{Kind: KindInt, Value: 42},
			{Kind: KindInt, Value: 42},
			{Kind: KindMapEnd},
		},
	},

	26: {
		value: map[any]any{
			42:    42,
			"foo": "bar",
		},
		expected: []Token{
			{Kind: KindMap},
			{Kind: KindString, Value: "foo"},
			{Kind: KindString, Value: "bar"},
			{Kind: KindInt, Value: 42},
			{Kind: KindInt, Value: 42},
			{Kind: KindMapEnd},
		},
	},

	27: {
		value: func() (int, string) {
			return 42, "42"
		},
		expected: []Token{
			{Kind: KindTuple},
			{Kind: KindInt, Value: 42},
			{Kind: KindString, Value: "42"},
			{Kind: KindTupleEnd},
		},
	},

	28: {
		value: marshalStringAsInt("foo"),
		expected: []Token{
			{Kind: KindInt, Value: int(3)},
		},
	},

	29: {
		value: testBool(true),
		expected: []Token{
			{Kind: KindBool, Value: true},
		},
	},

	30: {
		value: testInt8(42),
		expected: []Token{
			{Kind: KindInt8, Value: int8(42)},
		},
	},

	31: {
		value: testInt16(42),
		expected: []Token{
			{Kind: KindInt16, Value: int16(42)},
		},
	},

	32: {
		value: testInt32(42),
		expected: []Token{
			{Kind: KindInt32, Value: int32(42)},
		},
	},

	33: {
		value: testInt64(42),
		expected: []Token{
			{Kind: KindInt64, Value: int64(42)},
		},
	},

	34: {
		value: testUint(42),
		expected: []Token{
			{Kind: KindUint, Value: uint(42)},
		},
	},

	35: {
		value: testUint16(42),
		expected: []Token{
			{Kind: KindUint16, Value: uint16(42)},
		},
	},

	36: {
		value: testUint64(42),
		expected: []Token{
			{Kind: KindUint64, Value: uint64(42)},
		},
	},

	37: {
		value: testFloat32(42),
		expected: []Token{
			{Kind: KindFloat32, Value: float32(42)},
		},
	},

	38: {
		value: testFloat64(42),
		expected: []Token{
			{Kind: KindFloat64, Value: float64(42)},
		},
	},

	39: {
		value: testString("42"),
		expected: []Token{
			{Kind: KindString, Value: "42"},
		},
	},

	40: {
		value: testFloat32(math.NaN()),
		expected: []Token{
			{Kind: KindNaN},
		},
	},

	41: {
		value: testFloat64(math.NaN()),
		expected: []Token{
			{Kind: KindNaN},
		},
	},

	42: {
		value: definedInt(32),
		expected: []Token{
			{Kind: KindTypeName, Value: "github.com/reusee/sb.definedInt"},
			{Kind: KindInt, Value: int(32)},
		},
	},

	//
}

func TestMarshaler(t *testing.T) {

	for i, c := range marshalTestCases {

		// marshal
		tokens, err := TokensFromStream(MarshalCtx(c.ctx, c.value))
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
			MarshalCtx(c.ctx, c.value),
			Encode(buf),
		); err != nil {
			t.Fatal(err)
		}
		decoder := Decode(buf)
		if MustCompare(decoder, Tokens(c.expected).Iter()) != 0 {
			t.Fatalf("%d fail %+v", i, c)
		}

		// compare
		tokens, err = TokensFromStream(MarshalCtx(c.ctx, c.value))
		if err != nil {
			t.Fatal(err)
		}
		var obj any
		if err := Copy(tokens.Iter(), Unmarshal(&obj)); err != nil {
			t.Fatal(err)
		}
		if MustCompare(MarshalCtx(c.ctx, obj), MarshalCtx(c.ctx, c.value)) != 0 {
			t.Fatalf("not equal, got %#v, expected %#v", obj, c.value)
		}

		// CollectValueTokens
		var valueTokens Tokens
		if err := Copy(
			tokens.Iter(),
			CollectValueTokens(&valueTokens),
		); err != nil {
			t.Fatal(err)
		}
		if len(valueTokens) != len(tokens) {
			t.Fatalf("%d: expected %d, got %d\n", i, len(tokens), len(valueTokens))
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

func (badBinaryMarshaler) MarshalBinary() ([]byte, error) {
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

func (badTextMarshaler) MarshalText() ([]byte, error) {
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
	return func(token *Token) (Proc, error) {
		return MarshalValue(ctx, reflect.ValueOf(len(m)), cont), nil
	}
}

func TestValueMarshalFunc(t *testing.T) {
	fn := func(ctx Ctx, value reflect.Value, cont Proc) Proc {
		return MarshalValue(ctx, value, cont)
	}
	for _, c := range marshalTestCases {
		ctx := c.ctx
		ctx.Marshal = fn
		proc := fn(ctx, reflect.ValueOf(c.value), nil)
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
	proc1 := MarshalValue(Ctx{}.SkipEmpty(), reflect.ValueOf(struct {
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

func TestMarshalPath(t *testing.T) {

	// array
	if err := Copy(
		TapMarshal(DefaultCtx, [3]int{1, 2, 3}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if _, ok := v.([3]int); ok {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if i, ok := v.(int); ok && i == 1 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if i == 2 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			} else if i == 3 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 2 {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

	// slice
	if err := Copy(
		TapMarshal(DefaultCtx, []int{2, 3, 4}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if _, ok := v.([]int); ok {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if i, ok := v.(int); ok && i == 2 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if i == 3 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			} else if i == 4 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 2 {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

	// struct
	if err := Copy(
		TapMarshal(DefaultCtx, struct {
			Foo string
			Bar int
		}{
			Foo: "42",
			Bar: 42,
		}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if value.Kind() == reflect.Struct {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "Foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "42" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "Bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			} else if s, ok := v.(int); ok && s == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

	// map
	if err := Copy(
		TapMarshal(DefaultCtx, map[any]any{
			"Foo": "foo",
			"Bar": 42,
		}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if value.Kind() == reflect.Map {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "Foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Foo" {
					t.Fatal()
				}
			} else if s, ok := v.(string); ok && s == "Bar" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			} else if v == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != "Bar" {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

	// tuple
	if err := Copy(
		TapMarshal(DefaultCtx, func() (int, string) {
			return 42, "foo"
		}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if value.Kind() == reflect.Func {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if v == 42 {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
			} else if v == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

	// nested
	if err := Copy(
		TapMarshal(DefaultCtx, func() ([]int, string) {
			return []int{1, 2, 3}, "foo"
		}, func(ctx Ctx, value reflect.Value) {
			v := value.Interface()
			if value.Kind() == reflect.Func {
				if len(ctx.Path) != 0 {
					t.Fatal()
				}
			} else if v == 2 {
				if len(ctx.Path) != 2 {
					t.Fatal()
				}
				if ctx.Path[0] != 0 {
					t.Fatal()
				}
				if ctx.Path[1] != 1 {
					t.Fatal()
				}
			} else if v == "foo" {
				if len(ctx.Path) != 1 {
					t.Fatal()
				}
				if ctx.Path[0] != 1 {
					t.Fatal()
				}
			}
		}),
		Discard,
	); err != nil {
		t.Fatal(err)
	}

}

func TestMarshalBadTuple(t *testing.T) {
	err := Copy(
		Marshal(func(int) {}),
		Discard,
	)
	if !is(err, MarshalError) {
		t.Fatal()
	}
	if !is(err, BadTupleType) {
		t.Fatal()
	}
}

func TestMarshalNilSBMarshaler(t *testing.T) {
	var tokens Tokens
	var p *Custom
	if err := Copy(
		Marshal(p),
		CollectTokens(&tokens),
	); err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 {
		t.Fatal()
	}
	if tokens[0].Kind != KindNil {
		t.Fatal()
	}
}

func TestNilTupleFunc(t *testing.T) {
	var f func() (int, string)
	err := Copy(
		Marshal(f),
		Discard,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIgnoreFuncs(t *testing.T) {
	proc := MarshalValue(Ctx{
		IgnoreFuncs: true,
	}, reflect.ValueOf(struct {
		F func(int)
		I int
	}{
		F: func(int) {},
		I: 42,
	}), nil)
	var data struct {
		I int
	}
	if err := Copy(
		&proc,
		Unmarshal(&data),
	); err != nil {
		t.Fatal(err)
	}
	if data.I != 42 {
		t.Fatal()
	}
}
