package sb

import "testing"

func TestCompare(t *testing.T) {
	type Case [2]any

	type Foo struct {
		I int
		S string
	}

	equalCases := []Case{
		{true, true},
		{42, 42},
		{uint32(42), uint32(42)},
		{42.0, 42.0},
		{"foo", "foo"},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{Foo{42, "foo"}, Foo{42, "foo"}},
		{map[int]int{1: 1}, map[int]int{1: 1}},
		{Min, Min},
		{Max, Max},
	}
	for _, c := range equalCases {
		if MustCompare(NewMarshaler(c[0]), NewMarshaler(c[1])) != 0 {
			t.Fatal()
		}
	}

	notEqualCases := []Case{
		{false, true},
		{41, 42},
		{uint8(41), uint32(42)},
		{41.0, 42.0},
		{"fo", "foo"},
		{[]int{1, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 2}, []int{1, 2, 3}},
		{Foo{41, "foo"}, Foo{42, "foo"}},
		{Foo{42, "aoo"}, Foo{42, "foo"}},
		{42, []int{1, 2, 3}},
		{
			func() {},
			42,
		},
		{int8(42), int8(84)},
		{int16(42), int16(84)},
		{int32(42), int32(84)},
		{int64(42), int64(84)},
		{uint(42), uint(84)},
		{uint8(42), uint8(84)},
		{uint16(42), uint16(84)},
		{uint32(42), uint32(84)},
		{uint64(42), uint64(84)},
		{float32(42), float32(84)},
		{float64(42), float64(84)},
		{map[int]int{1: 1}, map[int]int{1: 42}},
		{[]byte("foo"), []byte("foobar")},
		{Min, 42},
		{42, Max},
	}
	for _, c := range notEqualCases {
		if MustCompare(NewMarshaler(c[0]), NewMarshaler(c[1])) != -1 {
			t.Fatal()
		}
		if MustCompare(NewMarshaler(c[1]), NewMarshaler(c[0])) != 1 {
			t.Fatal()
		}
	}

	if MustCompare(Tokens{
		{Kind: KindMax},
	}.Iter(), Tokens{
		{KindInt, 42},
	}.Iter()) != 1 {
		t.Fatal()
	}
	if MustCompare(Tokens{
		{Kind: KindMax},
	}.Iter(), Tokens{
		{Kind: KindMax},
	}.Iter()) != 0 {
		t.Fatal()
	}

}

func TestCompareBadStream(t *testing.T) {
	s1 := NewMarshaler(new(badTextMarshaler))
	s2 := NewMarshaler(42)
	_, err := Compare(s1, s2)
	if err == nil {
		t.Fatal()
	}
	_, err = Compare(s2, s1)
	if err == nil {
		t.Fatal()
	}
	func() {
		defer func() {
			if p := recover(); p == nil {
				t.Fatal()
			}
		}()
		MustCompare(s1, s2)
	}()
}
