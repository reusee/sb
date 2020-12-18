package sb

import (
	"testing"

	"github.com/reusee/e4"
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
	var stack *e4.Stacktrace
	if !as(err, &stack) {
		t.Fatal()
	}
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
