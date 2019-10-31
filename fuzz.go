package sb

import (
	"bytes"
	"reflect"
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
		panic(err)
	}
	bs := buf.Bytes()

	res, err := Compare(NewMarshaler(v), NewDecoder(bytes.NewReader(bs)))
	if err != nil {
		panic(err)
	}
	if res != 0 {
		pt("%d\n", res)
		pt("%+v\n", MustTokensFromStream(NewMarshaler(v)))
		pt("%+v\n", MustTokensFromStream(NewDecoder(bytes.NewReader(bs))))
		pt("%#v\n", v)
		panic(err)
	}

	var v2 any
	if err := Unmarshal(NewDecoder(bytes.NewReader(bs)), &v2); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(v, v2) {
		pt("%+v\n", MustTokensFromStream(NewMarshaler(v)))
		pt("%+v\n", MustTokensFromStream(NewDecoder(bytes.NewReader(bs))))
		pt("%#v\n", v)
		panic(err)
	}

	return 1
}
