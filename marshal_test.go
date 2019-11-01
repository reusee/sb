package sb

import (
	"bytes"
	"encoding"
	"fmt"
	"strings"
	"testing"
	"time"
)

type MarshalTestCase struct {
	value    any
	expected []Token
}

type foo int

var marshalTestCases = []MarshalTestCase{
	{
		int(42),
		[]Token{
			{KindInt, int(42)},
		},
	},

	{
		func() *int32 {
			i := int32(42)
			return &i
		}(),
		[]Token{
			{KindInt32, int32(42)},
		},
	},

	{
		true,
		[]Token{
			{KindBool, true},
		},
	},

	{
		uint32(42),
		[]Token{
			{KindUint32, uint32(42)},
		},
	},

	{
		float32(42),
		[]Token{
			{KindFloat32, float32(42)},
		},
	},

	{
		[]int{42, 4, 2},
		[]Token{
			{Kind: KindArray},
			{KindInt, int(42)},
			{KindInt, int(4)},
			{KindInt, int(2)},
			{Kind: KindArrayEnd},
		},
	},

	{
		[][]int{
			{42, 4, 2},
			{2, 4, 42},
		},
		[]Token{
			{Kind: KindArray},
			{Kind: KindArray},
			{KindInt, int(42)},
			{KindInt, int(4)},
			{KindInt, int(2)},
			{Kind: KindArrayEnd},
			{Kind: KindArray},
			{KindInt, int(2)},
			{KindInt, int(4)},
			{KindInt, int(42)},
			{Kind: KindArrayEnd},
			{Kind: KindArrayEnd},
		},
	},

	{
		"foo",
		[]Token{
			{KindString, "foo"},
		},
	},

	{
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
			{KindString, "Foo"},
			{KindInt, int(42)},
			{KindString, "Bar"},
			{KindFloat32, float32(42)},
			{KindString, "Baz"},
			{KindString, "42"},
			{KindString, "Boo"},
			{KindBool, false},
			{Kind: KindObjectEnd},
		},
	},

	{
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
			{KindInt, int(42)},
			{Kind: KindNil},
			{Kind: KindNil},
			{Kind: KindArrayEnd},
		},
	},

	{
		func() *int {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *bool {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *uint8 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *float32 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *[]int32 {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *[]int32 {
			array := []int32{42}
			return &array
		}(),
		[]Token{
			{Kind: KindArray},
			{KindInt32, int32(42)},
			{Kind: KindArrayEnd},
		},
	},

	{
		func() *string {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	{
		func() *struct{} {
			return nil
		}(),
		[]Token{
			{Kind: KindNil},
		},
	},

	func() MarshalTestCase {
		str := strings.Repeat("foo", 1024)
		return MarshalTestCase{
			str,
			[]Token{
				{KindString, str},
			},
		}
	}(),

	{
		foo(42),
		[]Token{
			{KindInt, int(42)},
		},
	},

	{
		[]foo{
			42,
		},
		[]Token{
			{Kind: KindArray},
			{KindInt, int(42)},
			{Kind: KindArrayEnd},
		},
	},

	{
		struct {
			Foo foo
		}{
			42,
		},
		[]Token{
			{Kind: KindObject},
			{KindString, "Foo"},
			{KindInt, int(42)},
			{Kind: KindObjectEnd},
		},
	},

	{
		struct {
			Foo []foo
		}{
			[]foo{
				42,
			},
		},
		[]Token{
			{Kind: KindObject},
			{KindString, "Foo"},
			{Kind: KindArray},
			{KindInt, int(42)},
			{Kind: KindArrayEnd},
			{Kind: KindObjectEnd},
		},
	},
}

func TestMarshaler(t *testing.T) {

	for i, c := range marshalTestCases {

		tokens, err := TokensFromStream(NewMarshaler(c.value))
		if err != nil {
			t.Fatal(err)
		}
		if len(tokens) != len(c.expected) {
			t.Fatalf("%d fail %+v", i, c)
		}
		for i, token := range tokens {
			if token != c.expected[i] {
				pt("expected %T\n", c.expected[i].Value)
				pt("token %T\n", token.Value)
				t.Fatalf("%d expected %#v, got %#v\nfail %+v", i, c.expected[i], token, c)
			}
		}

		buf := new(bytes.Buffer)
		if err := Encode(buf, NewMarshaler(c.value)); err != nil {
			t.Fatal(err)
		}
		decoder := NewDecoder(buf)
		if MustCompare(decoder, Tokens(c.expected).Iter()) != 0 {
			t.Fatalf("%d fail %+v", i, c)
		}

		tokens, err = TokensFromStream(NewMarshaler(c.value))
		if err != nil {
			t.Fatal(err)
		}
		var obj any
		if err := Unmarshal(tokens.Iter(), &obj); err != nil {
			t.Fatal(err)
		}
		if MustCompare(NewMarshaler(obj), NewMarshaler(c.value)) != 0 {
			t.Fatalf("not equal, got %#v, expected %#v", obj, c.value)
		}

	}

}

type Custom struct {
	Foo int
}

var _ Tokenizer = Custom{}

var _ Detokenizer = new(Custom)

func (c Custom) TokenizeSB() []Token {
	return []Token{
		{KindInt, c.Foo},
	}
}

func (c *Custom) DetokenizeSB(stream Stream) (err error) {
	p, err := stream.Next()
	if err != nil {
		return err
	}
	if p == nil {
		return
	}
	token := *p
	if token.Kind != KindInt {
		return
	}
	c.Foo = token.Value.(int)
	return
}

func TestCustomType(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Encode(buf, NewMarshaler(Custom{42})); err != nil {
		t.Fatal(err)
	}
	var c Custom
	if err := Unmarshal(NewDecoder(buf), &c); err != nil {
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
	m := NewMarshaler(v)
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
	m := NewMarshaler(v)
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
	m := NewMarshaler(now)
	tokens, err := TokensFromStream(m)
	if err != nil {
		t.Fatal(err)
	}
	var tt timeTextMarshaler
	if err := Unmarshal(tokens.Iter(), &tt); err != nil {
		t.Fatal(err)
	}
	if time.Since(tt.t) > time.Second {
		t.Fatal()
	}
}
