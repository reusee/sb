package sb

import "testing"

func TestTree(t *testing.T) {
	for _, c := range marshalTestCases {
		tokens, err := TokensFromStream(NewMarshaler(c.value))
		if err != nil {
			t.Fatal(err)
		}
		tree, err := TreeFromStream(tokens.Iter())
		if err != nil {
			t.Fatal(err)
		}
		res, err := Compare(tokens.Iter(), tree.Iter())
		if err != nil {
			t.Fatal(err)
		}
		if res != 0 {
			t.Fatal("not equal")
		}
	}
}
