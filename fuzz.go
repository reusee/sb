package sb

import (
	"bytes"
	"io/ioutil"
)

func Fuzz(data []byte) int {
	var v any
	err := Unmarshal(NewDecoder(bytes.NewReader(data)), &v)
	if err != nil {
		return 0
	}
	err = Encode(ioutil.Discard, NewMarshaler(v))
	if err != nil {
		return 0
	}
	return 1
}
