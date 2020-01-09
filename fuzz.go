package sb

import (
	"bytes"
	"crypto/md5"
	"hash/fnv"
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

		// hash
		hasher := NewPostHasher(NewMarshaler(obj2), md5.New)
		tokens, err := TokensFromStream(hasher)
		if err != nil {
			panic(err)
		}
		if tokens[len(tokens)-1].Kind != KindPostHash {
			panic("expecting hash token")
		}

		// sum
		sum1, err := HashSum(NewMarshaler(obj2), fnv.New128)
		if err != nil {
			panic(err)
		}
		sum2, err := HashSum(NewMarshaler(obj2), fnv.New128a)
		if err != nil {
			panic(err)
		}
		if bytes.Equal(sum1, sum2) {
			panic("should not equal")
		}

	}

	return 1 // NOCOVER
}
