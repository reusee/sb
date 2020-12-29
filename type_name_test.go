package sb

import (
	"reflect"
	"testing"
)

func TestTypeName(t *testing.T) {
	if TypeName(reflect.TypeOf((*string)(nil)).Elem()) != "string" {
		t.Fatal()
	}
	if TypeName(reflect.TypeOf((**string)(nil)).Elem()) != "*string" {
		t.Fatal()
	}
	if TypeName(reflect.TypeOf((***string)(nil)).Elem()) != "**string" {
		t.Fatal()
	}
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
}
