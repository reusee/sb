package sb

import (
	"reflect"
	"sync"
)

type HasDeprecatedFields interface {
	SBDeprecatedFields() []string
}

var hasDeprecatedFieldsType = reflect.TypeOf((*HasDeprecatedFields)(nil)).Elem()

var fieldIsDeprecatedMap sync.Map

type typeFieldName struct {
	Type reflect.Type
	Name string
}

func fieldIsDeprecated(
	t reflect.Type,
	fieldName string,
) (
	ret bool,
) {
	key := typeFieldName{
		Type: t,
		Name: fieldName,
	}
	if v, ok := fieldIsDeprecatedMap.Load(key); ok {
		return v.(bool)
	}
	defer func() {
		fieldIsDeprecatedMap.Store(key, ret)
	}()
	if !t.Implements(hasDeprecatedFieldsType) {
		return false
	}
	names := reflect.New(t).Elem().Interface().(HasDeprecatedFields).SBDeprecatedFields()
	for _, name := range names {
		if name == fieldName {
			return true
		}
	}
	return false
}
