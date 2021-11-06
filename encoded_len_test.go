package sb

import "testing"

func BenchmarkEncodedLen(b *testing.B) {
	obj := struct {
		Foo int
		Bar []int
		Baz map[int]func() (int, string)
	}{
		Bar: []int{1, 2, 3},
		Baz: map[int]func() (int, string){
			42: func() (int, string) {
				return 42, "foo"
			},
		},
	}
	var l int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Copy(
			Marshal(obj),
			EncodedLen(&l, nil),
		); err != nil {
			b.Fatal(err)
		}
	}
}
