package sb

import "testing"

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMarshaler(42)
		for {
			token, err := m.Next()
			if err != nil {
				b.Fatal(err)
			}
			if token == nil {
				break
			}
		}
	}
}
