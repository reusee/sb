package sb

import "testing"

func BenchmarkUnmarshalInt(b *testing.B) {
	n := int(42)
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(n),
			Unmarshal(&n),
		); err != nil {
			b.Fatal(err)
		}
	}
}
