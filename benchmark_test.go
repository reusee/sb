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

func BenchmarkUnmarshal(b *testing.B) {
	tokens := MustTokensFromStream(NewMarshaler(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var n int
		if err := Unmarshal(tokens.Iter(), &n); err != nil {
			b.Fatal(err)
		}
		if n != 42 {
			b.Fatal()
		}
	}
}
