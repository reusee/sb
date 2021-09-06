package sb

import (
	"testing"
)

func TestTree(t *testing.T) {
	for _, c := range marshalTestCases {
		tokens, err := TokensFromProc(MarshalCtx(c.ctx, c.value))
		if err != nil {
			t.Fatal(err)
		}
		tree, err := TreeFromProc(tokens.Iter())
		if err != nil {
			t.Fatal(err)
		}
		res, err := Compare(tokens.Iter(), tree.Iter())
		if err != nil {
			t.Fatal(err)
		}
		if res != 0 {
			t.Fatal("not equal")
		}
	}
}

func TestMoreThanOneValue(t *testing.T) {
	str := Tokens{
		{Kind: KindInt, Value: 42},
		{Kind: KindInt, Value: 42},
	}.Iter()
	_, err := TreeFromProc(str)
	if !is(err, MoreThanOneValue) {
		t.Fatal()
	}
}

func TestBadTreeFromProc(t *testing.T) {
	_, err := TreeFromProc(Tokens{
		{
			Kind: KindArrayEnd,
		},
	}.Iter())
	if !is(err, UnexpectedEndToken) {
		t.Fatal()
	}
}
