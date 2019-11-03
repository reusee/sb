package sb

import (
	"bytes"
	"io"
)

func Fuzz(data []byte) int { // NOCOVER
	r := bytes.NewReader(data)
	for {

		// decode and unmarshal
		if r.Len() == 0 {
			break
		}
		var obj any
		tee := new(bytes.Buffer)
		if err := Unmarshal(
			NewDecoder(io.TeeReader(r, tee)),
			&obj,
		); err != nil {
			return 0
		}
		teeBytes := tee.Bytes()

		// marshal and encode
		buf := new(bytes.Buffer)
		if err := Encode(buf, NewMarshaler(obj)); err != nil {
			panic(err)
		}
		bs := buf.Bytes()

		// decode and unmarshal
		var obj2 any
		if err := Unmarshal(
			NewDecoder(bytes.NewReader(teeBytes)),
			&obj2,
		); err != nil {
			panic(err)
		}

		// compare
		if MustCompare(
			NewDecoder(bytes.NewReader(bs)),
			NewMarshaler(obj2),
		) != 0 {
			tokens1 := MustTokensFromStream(
				NewDecoder(bytes.NewReader(teeBytes)),
			)
			tokens2 := MustTokensFromStream(
				NewDecoder(bytes.NewReader(bs)),
			)
			for i, token := range tokens1 {
				if i < len(tokens2) {
					pt("%+v\n%+v\n\n", token, tokens2[i])
				}
			}
			panic("not equal")
		}

	}

	return 1 // NOCOVER
}
