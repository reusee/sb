package sb

import (
	"testing"
)

func TestAltSink(t *testing.T) {
	var i int
	var b bool

	if err := Copy(
		Marshal(42),
		AltSink(
			Unmarshal(&i),
			Unmarshal(&b),
		),
	); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal()
	}
	if err := Copy(
		Marshal(24),
		AltSink(
			Unmarshal(&b),
			Unmarshal(&i),
		),
	); err != nil {
		t.Fatal(err)
	}
	if i != 24 {
		t.Fatal()
	}

	if err := Copy(
		Marshal(true),
		AltSink(
			Unmarshal(&i),
			Unmarshal(&b),
		),
	); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal()
	}
	if err := Copy(
		Marshal(false),
		AltSink(
			Unmarshal(&b),
			Unmarshal(&i),
		),
	); err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal()
	}

	err := Copy(
		Marshal("foo"),
		AltSink(
			Unmarshal(&i),
			Unmarshal(&b),
		),
	)
	if err == nil {
		t.Fatal()
	}

	var s string
	if err := Copy(
		Marshal("foo"),
		AltSink(
			Unmarshal(&b),
			Unmarshal(&i),
			Unmarshal(&s),
		),
	); err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal(err)
	}

	var ss []string
	if err := Copy(
		Marshal(
			[]string{"foo", "bar"},
		),
		AltSink(
			Unmarshal(&b),
			Unmarshal(&i),
			Unmarshal(&s),
			Unmarshal(&ss),
		),
	); err != nil {
		t.Fatal(err)
	}
	if len(ss) != 2 {
		t.Fatal()
	}

	var s1 struct {
		I int
	}
	var s2 struct{}
	if err := Copy(
		Marshal(
			struct {
				I int
			}{42},
		),
		AltSink(
			Unmarshal(&s1),
			Unmarshal(&s2),
		),
	); err != nil {
		t.Fatal(err)
	}
	if s1.I != 42 {
		t.Fatal()
	}

}

func TestExpectKind(t *testing.T) {
	err := Copy(
		Tokens{
			{
				Kind: KindInt,
			},
		}.Iter(),
		ExpectKind(KindString, nil),
	)
	if !is(err, ExpectingString) {
		t.Fatal()
	}

	err = Copy(
		Tokens{
			{
				Kind: KindInt,
			},
		}.Iter(),
		ExpectKind(KindInvalid, nil),
	)
	if !is(err, ExpectingValue) {
		t.Fatal()
	}
}

func TestFilterSink(t *testing.T) {
	var tokens Tokens
	if err := Copy(
		Marshal([]any{
			1, 2, 3, true, false, "foo", "bar",
		}),
		FilterSink(
			CollectTokens(&tokens),
			func(token *Token) bool {
				return token.Kind == KindBool
			},
		),
	); err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 7 {
		t.Fatal()
	}
}
