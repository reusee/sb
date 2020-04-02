package sb

import "testing"

func TestCopyNoExtraNext(t *testing.T) {
	var provide Proc
	tokens := Tokens{
		{
			Kind:  KindInt,
			Value: 1,
		},
		{
			Kind:  KindInt,
			Value: 2,
		},
		{
			Kind:  KindInt,
			Value: 3,
		},
	}
	provide = func() (*Token, Proc, error) {
		if len(tokens) == 0 {
			return nil, nil, nil
		}
		t := tokens[0]
		tokens = tokens[1:]
		return &t, provide, nil
	}
	var i int
	if err := Copy(
		&provide,
		Unmarshal(&i),
	); err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatal()
	}
	if len(tokens) != 2 {
		t.Fatal()
	}
}

func TestCopySingleShotSink(t *testing.T) {
	var ones Proc
	ones = func() (*Token, Proc, error) {
		return &Token{
			Kind:  KindInt,
			Value: 1,
		}, ones, nil
	}
	if err := Copy(
		&ones,
		nil,
	); err != nil {
		t.Fatal()
	}
}
