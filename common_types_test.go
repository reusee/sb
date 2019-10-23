package sb

import (
	"bytes"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	buf := new(bytes.Buffer)
	now := time.Now()
	if err := Encode(buf, NewMarshaler(now)); err != nil {
		t.Fatal(err)
	}
	var tt time.Time
	if err := Unmarshal(NewDecoder(buf), &tt); err != nil {
		t.Fatal(err)
	}
	if !tt.Equal(now) {
		t.Fatal()
	}
}
