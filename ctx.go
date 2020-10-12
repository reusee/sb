package sb

import "reflect"

type Ctx struct {
	Marshal   func(Ctx, reflect.Value, Proc) Proc
	Unmarshal func(Ctx, reflect.Value, Sink) Sink

	SkipEmptyStructFields       bool
	DisallowUnknownStructFields bool

	// array index, slice index, struct field name, map key, tuple index
	Path []any

	pointerDepth       int
	detectCycleEnabled bool
	visitedPointers    []uintptr
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
