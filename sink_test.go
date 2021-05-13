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

func TestSinkMarshal(t *testing.T) {
	var tokens Tokens
	sink := CollectTokens(&tokens)
	if _, err := sink.Marshal(42); err != nil {
		t.Fatal()
	}
	if len(tokens) != 1 {
		t.Fatal()
	}
}
