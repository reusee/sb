package sb

import (
	"bytes"
	"fmt"
	"hash"
)

var (
	NotFound = fmt.Errorf("not found")
)

func FindByHash(
	stream Stream,
	hash []byte,
	newState func() hash.Hash,
) (
	subStream Stream,
	err error,
) {

	var result *Tree

	_, err = TreeFromStream(
		stream,
		WithHash{newState},
		TapTree{
			func(tree *Tree) {
				if len(tree.Hash) > 0 &&
					bytes.Equal(tree.Hash, hash) {
					result = tree
				}
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, NotFound
	}

	return result.Iter(), nil
}
