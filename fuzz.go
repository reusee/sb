package sb

import (
	"bytes"
	"hash/fnv"
	"io"
	"math/rand"
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

		// validate tree
		_, err := TreeFromStream(NewDecoder(bytes.NewReader(teeBytes)))
		if err != nil {
			return 0
		}

		// marshal and encode
		buf := new(bytes.Buffer)
		if err := Encode(buf, NewMarshaler(obj)); err != nil { // NOCOVER
			panic(err)
		}
		bs := buf.Bytes()

		// decode and unmarshal
		var obj2 any
		if err := Unmarshal(
			NewDecoder(bytes.NewReader(teeBytes)),
			&obj2,
		); err != nil { // NOCOVER
			panic(err)
		}

		// compare
		if MustCompare(
			NewDecoder(bytes.NewReader(bs)),
			NewMarshaler(obj2),
		) != 0 { // NOCOVER
			tokens1 := MustTokensFromStream(
				NewDecoder(bytes.NewReader(teeBytes)),
			)
			tokens2 := MustTokensFromStream(
				NewDecoder(bytes.NewReader(bs)),
			)
			for i, token := range tokens1 { // NOCOVER
				if i < len(tokens2) { // NOCOVER
					pt("%+v\n%+v\n\n", token, tokens2[i])
				}
			}
			panic("not equal") // NOCOVER
		}

		// hash
		hasher := NewPostHasher(NewMarshaler(obj2), newMapHashState)
		hashedTokens, err := TokensFromStream(hasher)
		if err != nil { // NOCOVER
			panic(err)
		}
		if hashedTokens[len(hashedTokens)-1].Kind != KindPostTag { // NOCOVER
			panic("expecting tag token")
		}

		// sum
		sum1, err := MustTreeFromStream(NewMarshaler(obj2)).HashSum(fnv.New128)
		if err != nil { // NOCOVER
			panic(err)
		}
		sum2, err := MustTreeFromStream(NewMarshaler(obj2)).HashSum(fnv.New128a)
		if err != nil { // NOCOVER
			panic(err)
		}
		if bytes.Equal(sum1, sum2) { // NOCOVER
			panic("should not equal")
		}

		// sink hash
		var sum3 []byte
		if err := Unmarshal(NewMarshaler(obj2), Hasher(fnv.New128, &sum3, nil)); err != nil {
			panic(err)
		}
		if !bytes.Equal(sum1, sum3) {
			panic("should equal")
		}

		// tree
		tree, err := TreeFromStream(NewDecoder(bytes.NewReader(bs)))
		if err != nil { // NOCOVER
			panic(err)
		}
		if MustCompare(
			NewMarshaler(obj2),
			tree.Iter(),
		) != 0 { // NOCOVER
			panic("not equal")
		}

		// hashed tree
		hashedTree, err := TreeFromStream(
			NewPostHasher(NewMarshaler(obj2), fnv.New128),
		)
		if err != nil { // NOCOVER
			panic(err)
		}
		if MustCompare(
			hashedTree.Iter(),
			tree.Iter(),
		) != 0 { // NOCOVER
			panic("not equal")
		}
		h, ok := hashedTree.Tags.Get("hash")
		if !ok { // NOCOVER
			panic("no hash")
		}
		if !bytes.Equal(h, sum1) { // NOCOVER
			panic("hash not match")
		}

		mapHashSum, err := MustTreeFromStream(
			NewMarshaler(obj2),
		).HashSum(
			newMapHashState,
		)
		if err != nil { // NOCOVER
			panic(err)
		}

		// random transform
		transforms := []func(Stream) Stream{

			// marshal and unmarshal
			func(in Stream) Stream {
				var v any
				if err := Unmarshal(in, &v); err != nil { // NOCOVER
					panic(err)
				}
				return NewMarshaler(v)
			},

			// encode and decode
			func(in Stream) Stream {
				buf := new(bytes.Buffer)
				if err := Encode(buf, in); err != nil { // NOCOVER
					panic(err)
				}
				return NewDecoder(buf)
			},

			// post hasher
			func(in Stream) Stream {
				return NewPostHasher(in, newMapHashState)
			},

			// tokens
			func(in Stream) Stream {
				return MustTokensFromStream(in).Iter()
			},

			// tree iter
			func(in Stream) Stream {
				return MustTreeFromStream(in).Iter()
			},

			// tree func iter
			func(in Stream) Stream {
				return MustTreeFromStream(in).IterFunc(func(*Tree) (*Token, error) {
					return nil, nil
				})
			},

			// ref and deref
			func(in Stream) Stream {
				type ref struct {
					Hash []byte
					Tree *Tree
				}
				var refs []ref
				refed := MustTreeFromStream(in).IterFunc(func(tree *Tree) (*Token, error) {
					h, ok := tree.Tags.Get("hash")
					if !ok {
						return nil, nil
					}
					if rand.Intn(2) != 0 {
						return nil, nil
					}
					refs = append(refs, ref{
						Hash: h,
						Tree: tree,
					})
					return &Token{
						Kind:  KindRef,
						Value: h,
					}, nil
				})
				return Deref(refed, func(hash []byte) (Stream, error) {
					for _, ref := range refs {
						if bytes.Equal(ref.Hash, hash) {
							return ref.Tree.Iter(), nil
						}
					}
					panic("bad ref") // NOCOVER
				})
			},

			// stream iter
			func(in Stream) Stream {
				proc := IterStream(in, nil)
				return &proc
			},

			// random filter post hash
			func(in Stream) Stream {
				return Filter(in, func(token *Token) bool {
					return token.Kind == KindPostTag &&
						rand.Intn(2) == 0
				})
			},

			// find
			func(in Stream) Stream {
				sub, err := FindByHash(in, mapHashSum, newMapHashState)
				if err != nil { // NOCOVER
					panic(err)
				}
				return sub
			},

			// unmarshal to multiple
			func(in Stream) Stream {
				var ts [3]any
				if err := Unmarshal(in, &ts[0], &ts[1], &ts[2]); err != nil {
					panic(err)
				}
				return NewMarshaler(ts[rand.Intn(len(ts))])
			},
		}
		fn := func(in Stream) Stream {
			return in
		}
		for _, i := range rand.Perm(len(transforms) * 8) {
			i := i % len(transforms)
			f := fn
			fn = func(in Stream) Stream {
				return f(transforms[i](in))
			}
		}
		s := fn(tree.Iter())
		if MustCompare(s, tree.Iter()) != 0 { // NOCOVER
			panic("not equal")
		}
		sum1, err = MustTreeFromStream(tree.Iter()).HashSum(newMapHashState)
		if err != nil { // NOCOVER
			panic(err)
		}
		sum2, err = MustTreeFromStream(fn(tree.Iter())).HashSum(newMapHashState)
		if err != nil { // NOCOVER
			panic(err)
		}
		if !bytes.Equal(sum1, sum2) { // NOCOVER
			panic("hash not equal")
		}

	}

	return 1 // NOCOVER
}
