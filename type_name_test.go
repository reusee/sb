package sb

import (
	"reflect"
	"testing"
)

func TestTypeName(t *testing.T) {
	type S string
	if TypeName(reflect.TypeOf((*S)(nil)).Elem()) != "github.com/reusee/sb.S" {
		t.Fatal()
	}
	if TypeName(reflect.TypeOf((**S)(nil)).Elem()) != "*github.com/reusee/sb.S" {
		t.Fatal()
	}
	if name := TypeName(reflect.TypeOf((***S)(nil)).Elem()); name != "**github.com/reusee/sb.S" {
		t.Fatalf("got %s", name)
	}

	if TypeName(reflect.TypeOf((*int)(nil)).Elem()) != "int" {
		t.Fatal()
	}
}

func TestAliasTypeName(t *testing.T) {
	type Foo int
	type Bar = Foo
	nameFoo := TypeName(reflect.TypeOf((*Foo)(nil)).Elem())
	nameBar := TypeName(reflect.TypeOf((*Bar)(nil)).Elem())
	if nameFoo != nameBar {
		t.Fatal()
	}
}
