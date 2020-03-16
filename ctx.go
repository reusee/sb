package sb

import "reflect"

type Ctx struct {
	Marshal   func(Ctx, reflect.Value, Proc) Proc
	Unmarshal func(Ctx, reflect.Value, Sink) Sink

	ReserveStructFieldsOrder    bool
	SkipEmptyStructFields       bool
	DisallowUnknownStructFields bool
}
