package sb

import (
	"reflect"
	"testing"
)

func TestAltSink(t *testing.T) {
	var i int
	var b bool

	if err := Copy(
		NewMarshaler(42),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&i), nil),
			UnmarshalValue(reflect.ValueOf(&b), nil),
		),
	); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal()
	}
	if err := Copy(
		NewMarshaler(24),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&b), nil),
			UnmarshalValue(reflect.ValueOf(&i), nil),
		),
	); err != nil {
		t.Fatal(err)
	}
	if i != 24 {
		t.Fatal()
	}

	if err := Copy(
		NewMarshaler(true),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&i), nil),
			UnmarshalValue(reflect.ValueOf(&b), nil),
		),
	); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal()
	}
	if err := Copy(
		NewMarshaler(false),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&b), nil),
			UnmarshalValue(reflect.ValueOf(&i), nil),
		),
	); err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal()
	}

	err := Copy(
		NewMarshaler("foo"),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&i), nil),
			UnmarshalValue(reflect.ValueOf(&b), nil),
		),
	)
	if err == nil {
		t.Fatal()
	}

	var s string
	if err := Copy(
		NewMarshaler("foo"),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&b), nil),
			UnmarshalValue(reflect.ValueOf(&i), nil),
			UnmarshalValue(reflect.ValueOf(&s), nil),
		),
	); err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal(err)
	}

	var ss []string
	if err := Copy(
		NewMarshaler(
			[]string{"foo", "bar"},
		),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&b), nil),
			UnmarshalValue(reflect.ValueOf(&i), nil),
			UnmarshalValue(reflect.ValueOf(&s), nil),
			UnmarshalValue(reflect.ValueOf(&ss), nil),
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
		NewMarshaler(
			struct {
				I int
			}{42},
		),
		AltSink(
			UnmarshalValue(reflect.ValueOf(&s1), nil),
			UnmarshalValue(reflect.ValueOf(&s2), nil),
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
