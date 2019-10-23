package sb

import (
	"bytes"
	"strings"
	"testing"
)

func TestMarshaler(t *testing.T) {
	type Case struct {
		value    any
		expected []Token
	}

	type foo int

	cases := []Case{

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

		func() Case {
			str := strings.Repeat("foo", 1024)
			return Case{
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

	for i, c := range cases {

		tokens := Tokens(c.value)
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

		tokens = tokens[:0]
		m := NewMarshaler(c.value)
		for token := m.Next(); token != nil; token = m.Next() {
			tokens = append(tokens, *token)
		}
		if len(tokens) != len(c.expected) {
			t.Fatalf("%d fail %+v", i, c)
		}
		for i, token := range tokens {
			if token != c.expected[i] {
				t.Fatalf("expected %#v, got %#v\nfail %+v", c.expected[i], token, c)
			}
		}

		buf := new(bytes.Buffer)
		if err := Encode(buf, NewMarshaler(c.value)); err != nil {
			t.Fatal(err)
		}
		decoder := NewDecoder(buf)
		l := List(c.expected)
		if Compare(decoder, l) != 0 {
			t.Fatalf("%d fail %+v", i, c)
		}

		tokens = Tokens(c.value)
		var obj any
		if err := Unmarshal(List(tokens), &obj); err != nil {
			t.Fatal(err)
		}
		if Compare(NewMarshaler(obj), NewMarshaler(c.value)) != 0 {
			t.Fatalf("not equal, got %#v, expected %#v", obj, c.value)
		}

	}

}

func TestBadType(t *testing.T) {
	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Fatal()
			}
		}()
		Tokens(func() {})
	}()
}
