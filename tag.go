package sb

import "bytes"

type Tags [][]byte

func (t Tags) Get(name string) ([]byte, bool) {
	prefix := []byte(name + ":")
	for _, tag := range t {
		if bytes.HasPrefix(tag, prefix) {
			return bytes.TrimPrefix(tag, prefix), true
		}
	}
	return nil, false
}

func (t *Tags) Set(name string, value []byte) {
	prefix := []byte(name + ":")
	for i, tag := range *t {
		if bytes.HasPrefix(tag, prefix) {
			(*t)[i] = append(prefix, value...)
			return
		}
	}
	*t = append(*t, append(prefix, value...))
}

func (t *Tags) Add(toAdd []byte) {
	for _, tag := range *t {
		if bytes.Equal(tag, toAdd) {
			return
		}
	}
	*t = append(*t, toAdd)
}

func (t Tags) Iter(index int, cont Proc) Proc {
	return func() (*Token, Proc, error) {
		if index >= len(t) {
			return nil, cont, nil
		}
		return &Token{
			Kind:  KindPostTag,
			Value: t[index],
		}, t.Iter(index+1, cont), nil
	}
}
