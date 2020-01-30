package sb

type Sink func(*Token) (Sink, error)
