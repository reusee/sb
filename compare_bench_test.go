package sb

import (
	"bytes"
	"testing"
)

func BenchmarkCompare(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(func() string {
			return "foo"
		}),
		Encode(buf),
	); err != nil {
		b.Fatal(err)
	}
	bs1 := buf.Bytes()

	buf = new(bytes.Buffer)
	if err := Copy(
		Marshal(func() string {
			return "FOOBAR"
		}),
		Encode(buf),
	); err != nil {
		b.Fatal(err)
	}
	bs2 := buf.Bytes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MustCompare(
			DecodeForCompare(bytes.NewReader(bs1)),
			DecodeForCompare(bytes.NewReader(bs2)),
		)
	}

}
