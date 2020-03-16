package sb

import "testing"

func TestSinkAsUnmarshaler(t *testing.T) {
	var a int
	if err := Copy(
		Marshal(42),
		Unmarshal(
			Unmarshal(
				Unmarshal(&a),
			),
		),
	); err != nil {
		t.Fatal(err)
	}
	if a != 42 {
		t.Fatal()
	}
}
