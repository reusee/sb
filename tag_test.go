package sb

import (
	"bytes"
	"testing"
)

func TestTag(t *testing.T) {
	var tags Tags
	_, ok := tags.Get("foo")
	if ok {
		t.Fatal()
	}
	tags.Set("foo", []byte("foo"))
	value, ok := tags.Get("foo")
	if !ok {
		t.Fatal()
	}
	if !bytes.Equal(value, []byte("foo")) {
		t.Fatal()
	}
	tags.Set("foo", []byte("foo"))
	value, ok = tags.Get("foo")
	if !ok {
		t.Fatal()
	}
	if !bytes.Equal(value, []byte("foo")) {
		t.Fatal()
	}
	tags.Add([]byte("foo:foo"))
	if len(tags) != 1 {
		t.Fatal()
	}
}
