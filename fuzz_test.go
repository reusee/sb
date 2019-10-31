package sb

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestFuzzCorpus(t *testing.T) {
	files, err := filepath.Glob("corpus/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		var v any
		err = Unmarshal(NewDecoder(bytes.NewReader(data)), &v)
		if err != nil {
			continue
		}

		buf := new(bytes.Buffer)
		err = Encode(buf, NewMarshaler(v))
		if err != nil {
			t.Fatal(err)
		}
		bs := buf.Bytes()

		res, err := Compare(NewMarshaler(v), NewDecoder(bytes.NewReader(bs)))
		if err != nil {
			t.Fatal(err)
		}
		if res != 0 {
			pt("%d\n", res)
			pt("%+v\n", MustTokensFromStream(NewMarshaler(v)))
			pt("%+v\n", MustTokensFromStream(NewDecoder(bytes.NewReader(bs))))
			pt("%#v\n", v)
			t.Fatal(err)
		}

	}

}
