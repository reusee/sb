package sb

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func TestSinkHash(t *testing.T) {
	type Case struct {
		value    any
		expected string
		sums     []string
		kinds    []Kind
	}
	cases := []Case{

		{
			42,
			"0fcc339bcc03b2d67d97d0e2fa60bd41",
			[]string{
				"",
				"0fcc339bcc03b2d67d97d0e2fa60bd41",
			},
			[]Kind{
				KindInt,
				KindInt,
			},
		},

		{
			[]int{1, 2, 3},
			"1686c4524aa5e66d9cf9b98296ea178c",
			[]string{
				"",
				"",
				"62aabcb77703b2d6d746d674187a9a50",
				"",
				"255cd04db403b2d6e416b2ad65ec0309",
				"",
				"8f217470f503b2d6dfd16944f6c63576",
				"",
				"d228cb69101a8caf78912b704e4a1475",
				"1686c4524aa5e66d9cf9b98296ea178c",
			},
			[]Kind{
				KindArray,
				KindInt,
				KindInt,
				KindInt,
				KindInt,
				KindInt,
				KindInt,
				KindArrayEnd,
				KindArrayEnd,
				KindArray,
			},
		},
	}

	for i, c := range cases {
		var sum []byte
		if err := Copy(
			Marshal(c.value),
			Hash(fnv.New128, &sum, nil),
		); err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sum) != c.expected {
			t.Fatalf("%d: %#v, got %x", i, c.value, sum)
		}

		var sums []string
		var kinds []Kind
		if err := Copy(
			Marshal(c.value),
			HashFunc(
				fnv.New128,
				&sum,
				func(s []byte, token *Token) error {
					sums = append(sums, fmt.Sprintf("%x", s))
					kinds = append(kinds, token.Kind)
					return nil
				},
				nil,
			),
		); err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sum) != c.expected {
			t.Fatalf("%d: %#v, got %x", i, c.value, sum)
		}

		if MustCompare(
			Marshal(sums),
			Marshal(c.sums),
		) != 0 {
			for _, s := range sums {
				pt("%s\n", s)
			}
			t.Fatal()
		}
		if MustCompare(
			Marshal(kinds),
			Marshal(c.kinds),
		) != 0 {
			for _, k := range kinds {
				pt("%s\n", k.String())
			}
			t.Fatal()
		}

	}

}

func TestBadHash(t *testing.T) {

	if err := Copy(
		Tokens{}.Iter(),
		Hash(fnv.New128, nil, nil),
	); !is(err, ExpectingValue) {
		t.Fatal()
	}

}
