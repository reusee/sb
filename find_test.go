package sb

import (
	"bytes"
	"testing"
)

func TestFindByHash(t *testing.T) {
	h, err := HashSum(NewMarshaler(42), newMapHashState)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range []any{
		42,
		[]int{42},
		map[int]int{
			0: 42,
		},
	} {

		sub, err := FindByHash(
			NewPostHasher(NewMarshaler(v), newMapHashState),
			h,
			newMapHashState,
		)
		if err != nil {
			t.Fatal(err)
		}
		subHash, err := HashSum(sub, newMapHashState)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(subHash, h) {
			t.Fatal("bad hash")
		}

		sub, err = FindByHash(
			NewMarshaler(v), // no post hash
			h,
			newMapHashState,
		)
		if err != nil {
			t.Fatal(err)
		}
		subHash, err = HashSum(sub, newMapHashState)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(subHash, h) {
			t.Fatal("bad hash")
		}

		_, err = FindByHash(
			NewMarshaler(v),
			[]byte{},
			newMapHashState,
		)
		if !is(err, NotFound) {
			t.Fatal()
		}

	}

}

func TestBadFind(t *testing.T) {
	_, err := FindByHash(
		Tokens{
			{
				Kind:  KindPostHash,
				Value: []byte("foo"),
			},
		}.Iter(),
		[]byte("foo"),
		newMapHashState,
	)
	if !is(err, UnexpectedHashToken) {
		t.Fatal()
	}

	_, err = FindByHash(
		Tokens{
			{
				Kind: KindArrayEnd,
			},
		}.Iter(),
		[]byte("foo"),
		newMapHashState,
	)
	if !is(err, UnexpectedEndToken) {
		t.Fatal()
	}

	_, err = FindByHash(
		Tokens{
			{
				Kind:  KindInt,
				Value: 42,
			},
			{
				Kind:  KindInt,
				Value: 42,
			},
		}.Iter(),
		[]byte("foo"),
		newMapHashState,
	)
	if !is(err, MoreThanOneValue) {
		t.Fatal()
	}
}
