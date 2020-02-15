package sb

type Stream interface {
	Next() (*Token, error)
}
