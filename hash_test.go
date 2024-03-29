package sb

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
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

func TestHash2(t *testing.T) {
	var sum1, sum2 []byte
	if err := Copy(
		Tokens{
			{
				Kind:  KindInt,
				Value: 42,
			},
		}.Iter(),
		Hash(fnv.New128, &sum1, nil),
	); err != nil {
		t.Fatal(err)
	}
	if err := Copy(
		Marshal(42),
		Hash(fnv.New128, &sum2, nil),
	); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sum1, sum2) {
		t.Fatal()
	}

	if err := Copy(
		Tokens{
			{
				Kind: KindArray,
			},
			{
				Kind:  KindInt,
				Value: 42,
			},
			{
				Kind: KindArrayEnd,
			},
		}.Iter(),
		Hash(fnv.New128, &sum1, nil),
	); err != nil {
		t.Fatal(err)
	}
	if err := Copy(
		Marshal([]int{42}),
		Hash(fnv.New128, &sum2, nil),
	); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sum1, sum2) {
		t.Fatal()
	}
}

func TestBadHash(t *testing.T) {

	if err := Copy(
		Tokens{}.Iter(),
		Hash(fnv.New128, nil, nil),
	); !is(err, io.ErrUnexpectedEOF) {
		t.Fatal()
	}

	if err := Copy(
		Tokens{
			{
				Kind: KindObject,
			},
		}.Iter(),
		Hash(fnv.New128, nil, nil),
	); !is(err, io.ErrUnexpectedEOF) {
		t.Fatal()
	}

	Foo := errors.New("foo")
	var sum []byte
	if err := Copy(
		Marshal(42),
		HashFunc(fnv.New128, &sum, func(hash []byte, token *Token) error {
			return Foo
		}, nil),
	); !is(err, Foo) {
		t.Fatal()
	}

	if err := Copy(
		Marshal(42),
		HashFunc(fnv.New128, &sum, func(hash []byte, _ *Token) error {
			if len(hash) > 0 {
				return Foo
			}
			return nil
		}, nil),
	); !is(err, Foo) {
		t.Fatal()
	}

	if err := Copy(
		Marshal([]int{1, 2, 3}),
		HashFunc(fnv.New128, &sum, func(hash []byte, token *Token) error {
			if len(hash) > 0 && token.Kind == KindArrayEnd {
				return Foo
			}
			return nil
		}, nil),
	); !is(err, Foo) {
		t.Fatal()
	}

	if err := Copy(
		Marshal([]int{1, 2, 3}),
		HashFunc(fnv.New128, &sum, func(hash []byte, token *Token) error {
			if len(hash) == 0 && token.Kind == KindArrayEnd {
				return Foo
			}
			return nil
		}, nil),
	); !is(err, Foo) {
		t.Fatal()
	}

	if err := Copy(
		Marshal([]int{1, 2, 3}),
		HashFunc(fnv.New128, &sum, func(hash []byte, token *Token) error {
			if len(hash) > 0 && token.Kind == KindArray {
				return Foo
			}
			return nil
		}, nil),
	); !is(err, Foo) {
		t.Fatal()
	}

}

func TestRefHash(t *testing.T) {
	a := []any{
		1, 2, 3,
	}
	var hash []byte
	if err := Copy(
		Marshal(2),
		Hash(sha256.New, &hash, nil),
	); err != nil {
		t.Fatal(err)
	}
	b := []any{
		1, Ref(hash), 3,
	}

	var hashA []byte
	if err := Copy(
		Marshal(a),
		Hash(sha256.New, &hashA, nil),
	); err != nil {
		t.Fatal(err)
	}

	var hashB []byte
	if err := Copy(
		Marshal(b),
		Hash(sha256.New, &hashB, nil),
	); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(hashA, hashB) {
		t.Fatal()
	}
}

func newSHA256() hash.Hash {
	return sha256.New()
}

func benchmarkHashSHA256(b *testing.B, size int) {
	data := bytes.Repeat([]byte("a"), size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(size))
		if err := Copy(
			Marshal(data),
			Hash(newSHA256, nil, nil),
		); err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkHashSHA256(b *testing.B) {
	for i := 8; i < 8*1024*1024; i *= 4 {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			benchmarkHashSHA256(b, i)
		})
	}
}
