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
		hasher := NewPostHasher(NewMarshaler(obj2), newMapHashState)
		hashedTokens, err := TokensFromStream(hasher)
		if err != nil {
			panic(err)
		}
		if hashedTokens[len(hashedTokens)-1].Kind != KindPostHash {
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

		// tree
		tree, err := TreeFromStream(NewDecoder(bytes.NewReader(bs)))
		if err != nil {
			panic(err)
		}
		if MustCompare(
			NewMarshaler(obj2),
			tree.Iter(),
		) != 0 {
			panic("not equal")
		}

		// hashed tree
		hashedTree, err := TreeFromStream(
			NewPostHasher(NewMarshaler(obj2), fnv.New128),
		)
		if err != nil {
			panic(err)
		}
		if MustCompare(
			hashedTree.Iter(),
			tree.Iter(),
		) != 0 {
			panic("not equal")
		}
		if !bytes.Equal(hashedTree.Hash, sum1) {
			panic("hash not match")
		}

		mapHashSum, err := HashSum(
			NewMarshaler(obj2),
			newMapHashState,
		)
		if err != nil {
			panic(err)
		}

		// random transform
		transforms := []func(Stream) Stream{

			// marshal and unmarshal
			func(in Stream) Stream {
				var v any
				if err := Unmarshal(in, &v); err != nil {
					panic(err)
				}
				return NewMarshaler(v)
			},

			// encode and decode
			func(in Stream) Stream {
				buf := new(bytes.Buffer)
				if err := Encode(buf, in); err != nil {
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
					if len(tree.Hash) == 0 {
						return nil, nil
					}
					if rand.Intn(2) != 0 {
						return nil, nil
					}
					refs = append(refs, ref{
						Hash: tree.Hash,
						Tree: tree,
					})
					return &Token{
						Kind:  KindRef,
						Value: tree.Hash,
					}, nil
				})
				return Deref(refed, func(hash []byte) (Stream, error) {
					for _, ref := range refs {
						if bytes.Equal(ref.Hash, hash) {
							return ref.Tree.Iter(), nil
						}
					}
					panic("bad ref")
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
					return token.Kind == KindPostHash &&
						rand.Intn(2) == 0
				})
			},

			// find
			func(in Stream) Stream {
				sub, err := FindByHash(in, mapHashSum, newMapHashState)
				if err != nil {
					panic(err)
				}
				return sub
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
		if MustCompare(s, tree.Iter()) != 0 {
			panic("not equal")
		}
		sum1, err = HashSum(tree.Iter(), newMapHashState)
		if err != nil {
			panic(err)
		}
		sum2, err = HashSum(fn(tree.Iter()), newMapHashState)
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(sum1, sum2) {
			panic("hash not equal")
		}

	}

	return 1 // NOCOVER
}
