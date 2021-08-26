package sb

import (
	"bytes"
	"testing"
)

func TestDecodeJson(t *testing.T) {
	var n float64
	if err := Copy(
		DecodeJson(bytes.NewReader([]byte(`42`)), nil),
		Unmarshal(&n),
	); err != nil {
		t.Fatal(err)
	}
	if n != 42 {
		t.Fatal()
	}

	var i int64
	if err := Copy(
		DecodeJson(bytes.NewReader([]byte(`42`)), nil),
		Unmarshal(&i),
	); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatalf("get %d\n", i)
	}

	var f float64
	if err := Copy(
		DecodeJson(bytes.NewReader([]byte(`4.2`)), nil),
		Unmarshal(&f),
	); err != nil {
		t.Fatal(err)
	}
	if f != 4.2 {
		t.Fatalf("get %v\n", f)
	}

}
