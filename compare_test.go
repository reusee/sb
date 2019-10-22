package sb

import "testing"

func TestCompare(t *testing.T) {
	type Case [2]any

	type Foo struct {
		I int
		S string
	}

	cases := []Case{
		{true, true},
		{42, 42},
		{uint32(42), uint8(42)},
		{42.0, 42.0},
		{"foo", "foo"},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{Foo{42, "foo"}, Foo{42, "foo"}},
	}
	for _, c := range cases {
		if Compare(c[0], c[1]) != 0 {
			t.Fatal()
		}
	}

	cases = []Case{
		{false, true},
		{41, 42},
		{uint32(41), uint8(42)},
		{41.0, 42.0},
		{"fo", "foo"},
		{[]int{1, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 2}, []int{1, 2, 3}},
		{Foo{41, "foo"}, Foo{42, "foo"}},
		{Foo{42, "aoo"}, Foo{42, "foo"}},
		{42, []int{1, 2, 3}},
	}
	for _, c := range cases {
		if Compare(c[0], c[1]) != -1 {
			t.Fatal()
		}
		if Compare(c[1], c[0]) != 1 {
			t.Fatal()
		}
	}

}
