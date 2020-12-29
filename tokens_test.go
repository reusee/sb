package sb

import (
	"bytes"
	"testing"
)

func TestTokens(t *testing.T) {
	// bad stream
	func() {
		defer func() {
			if p := recover(); p == nil {
				t.Fatal()
			}
		}()
		MustTokensFromStream(Decode(bytes.NewReader([]byte{
			byte(KindString), // incomplete
		})))
	}()

	for _, c := range marshalTestCases {
		MustTokensFromStream(MarshalCtx(c.ctx, c.value))
	}

}

func TestCollectValueTokens(t *testing.T) {
	type Foo struct {
		I int
		S string
	}
	var tokens Tokens
	v := struct {
		I int
		S Sink
	}{
		S: CollectValueTokens(&tokens),
	}
	if err := Copy(
		Marshal(Foo{
			I: 42,
			S: "foo",
		}),
		Unmarshal(&v),
	); err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 {
		t.Fatal()
	}
	if tokens[0].Kind != KindString || tokens[0].Value != "foo" {
		t.Fatal()
	}

	err := Copy(
		Tokens{
			{Kind: KindArrayEnd},
		}.Iter(),
		Unmarshal(
			CollectValueTokens(&tokens),
		),
	)
	if !is(err, UnexpectedEndToken) {
		t.Fatal()
	}

	err = Copy(
		Tokens{
			{Kind: KindArray},
		}.Iter(),
		Unmarshal(
			CollectValueTokens(&tokens),
		),
	)
	if !is(err, ExpectingValue) {
		t.Fatal()
	}

	err = Copy(
		Tokens{
			{Kind: KindInt},
		}.Iter(),
		Unmarshal(
			CollectValueTokens(&tokens),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = Copy(
		Tokens{
			{Kind: KindArray},
			{Kind: KindObjectEnd},
		}.Iter(),
		Unmarshal(
			CollectValueTokens(&tokens),
		),
	)
	if !is(err, ExpectingArrayEnd) {
		t.Fatal()
	}

	err = Copy(
		Tokens{
			{Kind: KindArray},
			{Kind: KindArray},
			{Kind: KindArrayEnd},
			{Kind: KindObject},
			{Kind: KindObjectEnd},
			{Kind: KindMap},
			{Kind: KindMapEnd},
			{Kind: KindTuple},
			{Kind: KindTupleEnd},
			{Kind: KindArrayEnd},
		}.Iter(),
		Unmarshal(
			CollectValueTokens(&tokens),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

}
