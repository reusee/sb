package sb

import (
	"bytes"
	"hash"
)

func FindByHash(
	stream Stream,
	hash []byte,
	newState func() hash.Hash,
) (
	subStream Stream,
	err error,
) {

	tree, err := TreeFromStream(stream)
	if err != nil {
		return nil, err
	}
	if err := tree.FillHash(newState); err != nil { // NOCOVER
		return nil, err
	}

	var iter func(*Tree) (Stream, error)
	iter = func(tree *Tree) (Stream, error) {
		if bytes.Equal(tree.Hash, hash) {
			return tree.Iter(), nil
		}
		for _, sub := range tree.Subs {
			if subStream, err := iter(sub); err != nil {
				return nil, err
			} else if subStream != nil {
				return subStream, nil
			}
		}
		return nil, nil
	}

	subStream, err = iter(tree)
	if subStream == nil {
		err = NotFound
	}

	return
}
