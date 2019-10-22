package sb

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	type Case struct {
		value    any
		expected []Token
	}

	cases := []Case{

		{
			int(42),
			[]Token{
				{KindInt, int64(42)},
			},
		},

		{
			func() *int32 {
				i := int32(42)
				return &i
			}(),
			[]Token{
				{KindIndirect, KindInt},
				{KindInt, int64(42)},
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
				{KindUint, uint64(42)},
			},
		},

		{
			float32(42),
			[]Token{
				{KindFloat, float64(42)},
			},
		},

		{
			[]int{42, 4, 2},
			[]Token{
				{Kind: KindArray},
				{KindInt, int64(42)},
				{KindInt, int64(4)},
				{KindInt, int64(2)},
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
				{KindInt, int64(42)},
				{KindInt, int64(4)},
				{KindInt, int64(2)},
				{Kind: KindArrayEnd},
				{Kind: KindArray},
				{KindInt, int64(2)},
				{KindInt, int64(4)},
				{KindInt, int64(42)},
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
			}{
				42,
				42,
				"42",
			},
			[]Token{
				{Kind: KindObject},
				{KindString, "Foo"},
				{KindInt, int64(42)},
				{KindString, "Bar"},
				{KindFloat, float64(42)},
				{KindString, "Baz"},
				{KindString, "42"},
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
				t.Fatalf("expected %#v, got %#v\nfail %+v", c.expected[i], token, c)
			}
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
