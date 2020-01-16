package sb

import "hash"

func HashSum(
	stream Stream,
	newState func() hash.Hash,
) (
	sum []byte,
	err error,
) {
	tree, err := TreeFromStream(stream)
	if err != nil { // NOCOVER
		return nil, err
	}
	if err := tree.FillHash(newState); err != nil { // NOCOVER
		return nil, err
	}
	return tree.Hash, nil
}
