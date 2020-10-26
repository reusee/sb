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
		if err := Copy(
			Decode(io.TeeReader(r, tee)),
			Unmarshal(&obj),
		); err != nil {
			var decodeErr DecodeError
			var unmarshalErr UnmarshalError
			if !as(err, &decodeErr) && !as(err, &unmarshalErr) { // NOCOVER
				panic("should be decode or unmarshal error")
			}
			return 0
		}
		teeBytes := tee.Bytes()

		// validate tree
		_, err := TreeFromStream(Decode(bytes.NewReader(teeBytes)))
		if err != nil { // NOCOVER
			return 0
		}

		// marshal and encode
		buf := new(bytes.Buffer)
		if err := Copy(Marshal(obj), Encode(buf)); err != nil { // NOCOVER
			panic(err)
		}
		bs := buf.Bytes()

		// decode and unmarshal
		var obj2 any
		if err := Copy(
			Decode(bytes.NewReader(teeBytes)),
			Unmarshal(&obj2),
		); err != nil { // NOCOVER
			panic(err)
		}

		// compare
		var tokens1 Tokens
		var tokens2 Tokens
		if MustCompare(
			Tee(
				Decode(bytes.NewReader(bs)),
				CollectTokens(&tokens1),
			),
			Tee(
				Marshal(obj2),
				CollectValueTokens(&tokens2),
			),
		) != 0 { // NOCOVER
			pt("obj : %+v\n\n", obj)
			pt("obj2: %+v\n\n", obj2)
			for i, token := range tokens1 { // NOCOVER
				if i < len(tokens2) { // NOCOVER
					pt("%+v\n%+v\n\n", token, tokens2[i])
				}
			}
			for _, token := range tokens2[len(tokens1):] { // NOCOVER
				pt("---\n%+v\n\n", token)
			}
			panic("not equal") // NOCOVER
		}

		// sum
		var sum1, sum2 []byte
		if err := Copy(
			Marshal(obj2),
			Hash(fnv.New128, &sum1, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		if err := Copy(
			Marshal(obj2),
			Hash(fnv.New128a, &sum2, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		if bytes.Equal(sum1, sum2) { // NOCOVER
			panic("should not equal")
		}

		// sink hash
		var sum3 []byte
		if err := Copy(
			Marshal(obj2),
			Hash(fnv.New128, &sum3, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		if !bytes.Equal(sum1, sum3) { // NOCOVER
			panic("should equal")
		}

		// tree
		tree, err := TreeFromStream(Decode(bytes.NewReader(bs)))
		if err != nil { // NOCOVER
			panic(err)
		}
		if MustCompare(
			Marshal(obj2),
			tree.Iter(),
		) != 0 { // NOCOVER
			panic("not equal")
		}

		var mapHashSum []byte
		if err := Copy(
			Marshal(obj2),
			Hash(newMapHashState, &mapHashSum, nil),
		); err != nil { // NOCOVER
			panic(err)
		}

		// random transform
		transforms := []func(Stream) Stream{

			// marshal and unmarshal
			func(in Stream) Stream {
				var v any
				if err := Copy(
					in,
					Unmarshal(&v),
				); err != nil { // NOCOVER
					panic(err)
				}
				return Marshal(v)
			},

			// encode and decode
			func(in Stream) Stream {
				buf := new(bytes.Buffer)
				if err := Copy(in, Encode(buf)); err != nil { // NOCOVER
					panic(err)
				}
				return Decode(buf)
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
				tree := MustTreeFromStream(in)
				if err := tree.FillHash(newMapHashState); err != nil { // NOCOVER
					panic(err)
				}
				refed := tree.IterFunc(func(tree *Tree) (*Token, error) {
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
					panic("bad ref") // NOCOVER
				})
			},

			// stream iter
			func(in Stream) Stream {
				proc := IterStream(in, nil)
				return &proc
			},

			// marshal stream iter
			func(in Stream) Stream {
				return Marshal(
					IterStream(in, nil),
				)
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
				if err := Copy(
					Tee(
						in,
						Unmarshal(&ts[0]),
						Unmarshal(&ts[1]),
						Unmarshal(&ts[2]),
					),
					Discard,
				); err != nil { // NOCOVER
					panic(err)
				}
				return Marshal(ts[rand.Intn(len(ts))])
			},

			// tee
			func(in Stream) Stream {
				return Tee(in)
			},

			// tee 2
			func(in Stream) Stream {
				buf := new(bytes.Buffer)
				if err := Copy(
					Tee(in, Encode(buf)),
					Discard,
				); err != nil { // NOCOVER
					panic(err)
				}
				return Decode(buf)
			},

			// collect tokens
			func(in Stream) Stream {
				var tokens Tokens
				if err := Copy(in, CollectTokens(&tokens)); err != nil { // NOCOVER
					panic(err)
				}
				return tokens.Iter()
			},

			// collect value tokens
			func(in Stream) Stream {
				var tokens Tokens
				if err := Copy(in, CollectValueTokens(&tokens)); err != nil { // NOCOVER
					panic(err)
				}
				return tokens.Iter()
			},

			// tuple
			func(in Stream) Stream {
				var v any
				if err := Copy(in, Unmarshal(&v)); err != nil { // NOCOVER
					panic(err)
				}
				var tuple Tuple
				if err := Copy(
					Marshal(Tuple{v}),
					Unmarshal(&tuple),
				); err != nil { // NOCOVER
					panic(err)
				}
				return Marshal(tuple[0])
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
		if MustCompare(fn(tree.Iter()), tree.Iter()) != 0 { // NOCOVER
			panic("not equal")
		}

		if err := Copy(
			tree.Iter(),
			Hash(newMapHashState, &sum1, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		if err := Copy(
			fn(tree.Iter()),
			Hash(newMapHashState, &sum2, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		if !bytes.Equal(sum1, sum2) { // NOCOVER
			panic("hash not equal")
		}

	}

	return 1 // NOCOVER
}
