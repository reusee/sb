package sb

import (
	"testing"
)

func TestError(t *testing.T) {
	var err error

	var v map[int]map[int]string
	err = Copy(
		Marshal(map[int]map[int]int{
			42: {
				43: 44,
			},
		}),
		Unmarshal(&v),
	)
	if !is(err, ExpectingInt) {
		t.Fatal()
	}
	var path Path
	if !as(err, &path) {
		t.Fatal()
	}
	if path.String() != "/42/43" {
		t.Fatal()
	}
	if !is(err, UnmarshalError) {
		t.Fatal()
	}

}
