package sb

import (
	"errors"
	"fmt"

	"github.com/reusee/e4"
	"github.com/reusee/pp"
)

type (
	any = interface{}
)

var (
	pt = fmt.Printf
	is = errors.Is
	as = errors.As

	we = e4.Wrap
	ce = e4.Check

	Copy        = pp.Copy[Token, Proc, Sink]
	ConcatSinks = pp.CatSink[Token, Sink]
	ConcatProcs = pp.CatSrc[Token, Proc]
)
