package sb

import (
	"reflect"
	"sync"
)

var typeToName, nameToType sync.Map

func TypeName(t reflect.Type) (name string) {
	if v, ok := typeToName.Load(t); ok {
		return v.(string)
	}

	defer func() {
		typeToName.Store(t, name)
		nameToType.Store(name, t)
	}()

	// pointer
	if t.Kind() == reflect.Ptr {
		return "*" + TypeName(t.Elem())
	}

	// defined types with package path
	if definedName := t.Name(); definedName != "" {
		if pkgPath := t.PkgPath(); pkgPath != "" {
			return pkgPath + "." + definedName
		}
	}

	// fallback
	return t.String()
}

func Register(t reflect.Type) {
	TypeName(t)
}
