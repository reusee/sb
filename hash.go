package sb

import "hash"

func HashSum(
	stream Stream,
	newState func() hash.Hash,
) (
	sum []byte,
	err error,
) {
	hasher := NewPostHasher(stream, newState)
	var token, last *Token
	for {
		token, err = hasher.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			if last.Kind != KindPostHash {
				panic("bad hasher")
			}
			sum = last.Value.([]byte)
			break
		}
		last = token
	}
	return
}
