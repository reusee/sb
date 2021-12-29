package sb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/reusee/e4"
)

type Ctx struct {
	Marshal   func(Ctx, reflect.Value, Proc) Proc
	Unmarshal func(Ctx, reflect.Value, Sink) Sink

	visitedPointers []uintptr

	// array index, slice index, struct field name, map key, tuple index
	Path Path

	pointerDepth int

	SkipEmptyStructFields       bool
	DisallowUnknownStructFields bool
	detectCycleEnabled          bool

	IgnoreFuncs bool
}

type Path []any

func (p Path) String() string {
	var b strings.Builder
	for _, elem := range p {
		b.WriteString(fmt.Sprintf("/%v", elem))
	}
	return b.String()
}

func WithPath(ctx Ctx) e4.WrapFunc {
	return func(prev error) error {
		return e4.Error{
			Err:  append(ctx.Path[:0:0], ctx.Path...),
			Prev: prev,
		}
	}
}

var _ error = Path{}

func (p Path) Error() string {
	return "path: " + p.String()
}

var DefaultCtx = Ctx{
	Marshal:   MarshalValue,
	Unmarshal: UnmarshalValue,
}

func (c Ctx) SkipEmpty() Ctx {
	c.SkipEmptyStructFields = true
	return c
}

func (c Ctx) Strict() Ctx {
	c.DisallowUnknownStructFields = true
	return c
}

func (c Ctx) WithPath(path any) Ctx {
	c.Path = append(c.Path, path)
	return c
}
