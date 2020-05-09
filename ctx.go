package sb

import "reflect"

type Ctx struct {
	Marshal   func(Ctx, reflect.Value, Proc) Proc
	Unmarshal func(Ctx, reflect.Value, Sink) Sink

	ReserveStructFieldsOrder    bool
	SkipEmptyStructFields       bool
	DisallowUnknownStructFields bool
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
