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
			if !is(err, DecodeError) && !is(err, UnmarshalError) { // NOCOVER
				panic("should be decode or unmarshal error")
			}
			return 0
		}
		teeBytes := tee.Bytes()

		// validate tree
		_, err := TreeFromProc(Decode(bytes.NewReader(teeBytes)))
		if err != nil { // NOCOVER
			return 0
		}

		// marshal and encode
		buf := new(bytes.Buffer)
		var l int
		if err := Copy(
			Marshal(obj),
			Encode(buf),
			EncodedLen(&l, nil),
		); err != nil { // NOCOVER
			panic(err)
		}
		bs := buf.Bytes()
		if l != len(bs) {
			panic("bad len")
		}

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

		// compare bytes
		if MustCompareBytes(bs, bs) != 0 {
			panic("not equal")
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
		tree, err := TreeFromProc(Decode(bytes.NewReader(bs)))
		if err != nil { // NOCOVER
			panic(err)
		}
		if MustCompare(
			Marshal(obj2),
			tree.Iter(),
		) != 0 { // NOCOVER
			panic("not equal")
		}

		// map hash sum
		var mapHashSum []byte
		if err := Copy(
			Marshal(obj2),
			Hash(newMapHashState, &mapHashSum, nil),
		); err != nil { // NOCOVER
			panic(err)
		}

		// tree with hash
		treeWithHash, err := TreeFromProc(
			Decode(bytes.NewReader(bs)),
			WithHash{newMapHashState},
		)
		if err != nil { // NOCOVER
			panic(err)
		}
		if !bytes.Equal(treeWithHash.Hash, mapHashSum) { // NOCOVER
			panic("should equal")
		}

		// random transform
		transforms := []func(Proc) Proc{

			// marshal and unmarshal
			func(in Proc) Proc {
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
			func(in Proc) Proc {
				buf := new(bytes.Buffer)
				if err := Copy(in, Encode(buf)); err != nil { // NOCOVER
					panic(err)
				}
				return Decode(buf)
			},

			// tokens
			func(in Proc) Proc {
				return MustTokensFromProc(in).Iter()
			},

			// tree iter
			func(in Proc) Proc {
				return MustTreeFromProc(in).Iter()
			},

			// tree func iter
			func(in Proc) Proc {
				return MustTreeFromProc(in).IterFunc(func(*Tree) (*Token, error) {
					return nil, nil
				})
			},

			// ref and deref
			func(in Proc) Proc {
				type ref struct {
					Tree *Tree
					Hash []byte
				}
				var refs []ref
				tree := MustTreeFromProc(in)
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
				return Deref(refed, func(hash []byte) (Proc, error) {
					for _, ref := range refs {
						if bytes.Equal(ref.Hash, hash) {
							return ref.Tree.Iter(), nil
						}
					}
					panic("bad ref") // NOCOVER
				})
			},

			// iter
			func(in Proc) Proc {
				return Iter(in, nil)
			},

			// marshal iter
			func(in Proc) Proc {
				return Marshal(
					Iter(in, nil),
				)
			},

			// find
			func(in Proc) Proc {
				sub, err := FindByHash(in, mapHashSum, newMapHashState)
				if err != nil { // NOCOVER
					panic(err)
				}
				return sub
			},

			// unmarshal to multiple
			func(in Proc) Proc {
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
			func(in Proc) Proc {
				return Tee(in)
			},

			// tee 2
			func(in Proc) Proc {
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
			func(in Proc) Proc {
				var tokens Tokens
				if err := Copy(in, CollectTokens(&tokens)); err != nil { // NOCOVER
					panic(err)
				}
				return tokens.Iter()
			},

			// collect value tokens
			func(in Proc) Proc {
				var tokens Tokens
				if err := Copy(in, CollectValueTokens(&tokens)); err != nil { // NOCOVER
					panic(err)
				}
				return tokens.Iter()
			},

			// sink marshal
			func(in Proc) Proc {
				var tokens Tokens
				if _, err := CollectTokens(&tokens).Marshal(in); err != nil {
					panic(err)
				}
				return tokens.Iter()
			},

			// tuple
			func(in Proc) Proc {
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

		fn := func(in Proc) Proc {
			return in
		}
		for _, i := range rand.Perm(len(transforms) * 8) {
			i := i % len(transforms)
			f := fn
			fn = func(in Proc) Proc {
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
