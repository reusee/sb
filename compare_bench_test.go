package sb

import (
	"bytes"
	"testing"
)

func BenchmarkCompare(b *testing.B) {
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(func() int {
			return 42
		}),
		Encode(buf),
	); err != nil {
		b.Fatal(err)
	}
	bs1 := buf.Bytes()

	buf = new(bytes.Buffer)
	if err := Copy(
		Marshal(func() int {
			return 99
		}),
		Encode(buf),
	); err != nil {
		b.Fatal(err)
	}
	bs2 := buf.Bytes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MustCompare(
			Decode(bytes.NewReader(bs1)),
			Decode(bytes.NewReader(bs2)),
		)
	}

}
