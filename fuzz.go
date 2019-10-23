package sb

import (
	"bytes"
)

func Fuzz(data []byte) int {
	var v any
	err := Unmarshal(NewDecoder(bytes.NewReader(data)), &v)
	if err != nil {
		return 0
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, NewMarshaler(v))
	if err != nil {
		return 0
	}
	bs := buf.Bytes()

	res, err := Compare(NewMarshaler(v), NewDecoder(bytes.NewReader(bs)))
	if err != nil {
		panic(err)
	}
	if res != 0 {
		pt("%+v\n", MustTokensFromStream(NewMarshaler(v)))
		pt("%+v\n", MustTokensFromStream(NewDecoder(bytes.NewReader(bs))))
		pt("%#v\n", v)
		panic(err)
	}

	return 1
}
