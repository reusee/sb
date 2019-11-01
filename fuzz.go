package sb

import (
	"bytes"
)

func Fuzz(data []byte) int { // NOCOVER
	var v any
	err := Unmarshal(NewDecoder(bytes.NewReader(data)), &v)
	if err != nil { // NOCOVER
		return 0
	}

	buf := new(bytes.Buffer) // NOCOVER
	err = Encode(buf, NewMarshaler(v))
	if err != nil { // NOCOVER
		panic(err)
	}
	bs := buf.Bytes() // NOCOVER

	res, err := Compare(NewMarshaler(v), NewDecoder(bytes.NewReader(bs)))
	if err != nil { // NOCOVER
		panic(err)
	}
	if res != 0 { // NOCOVER
		pt("%d\n", res)
		pt("%+v\n", MustTokensFromStream(NewMarshaler(v)))
		pt("%+v\n", MustTokensFromStream(NewDecoder(bytes.NewReader(bs))))
		pt("%#v\n", v)
		panic(err)
	}

	return 1 // NOCOVER
}
