package sb

import (
	"fmt"
	"reflect"
	"strings"
)

type Ctx struct {
	Marshal   func(Ctx, reflect.Value, Proc) Proc
	Unmarshal func(Ctx, reflect.Value, Sink) Sink

	SkipEmptyStructFields       bool
	DisallowUnknownStructFields bool

	detectCycleEnabled bool
	pointerDepth       int
	visitedPointers    []uintptr

	// array index, slice index, struct field name, map key, tuple index
	Path Path
}

type Path []any

func (p Path) String() string {
	var b strings.Builder
	for _, elem := range p {
		b.WriteString(fmt.Sprintf("/%v", elem))
	}
	return b.String()
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
