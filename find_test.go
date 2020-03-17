package sb

import (
	"bytes"
	"testing"
)

func TestFindByHash(t *testing.T) {
	var h []byte
	err := Copy(
		Marshal(42),
		Hash(newMapHashState, &h, nil),
	)
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
			PostHash(Marshal(v), newMapHashState),
			h,
			newMapHashState,
		)
		if err != nil {
			t.Fatal(err)
		}
		var subHash []byte
		err = Copy(
			sub,
			Hash(newMapHashState, &subHash, nil),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(subHash, h) {
			t.Fatal("bad hash")
		}

		sub, err = FindByHash(
			Marshal(v), // no post hash
			h,
			newMapHashState,
		)
		if err != nil {
			t.Fatal(err)
		}
		err = Copy(
			sub,
			Hash(newMapHashState, &subHash, nil),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(subHash, h) {
			t.Fatal("bad hash")
		}

		_, err = FindByHash(
			Marshal(v),
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
				Kind:  KindPostTag,
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
