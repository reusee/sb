package sb

import (
	"fmt"
	"reflect"
	"sync"
)

var typeToName, registeredNameToType, registeredTypeToName sync.Map

func TypeName(t reflect.Type) (name string) {
	if v, ok := typeToName.Load(t); ok {
		return v.(string)
	}

	defer func() {
		typeToName.Store(t, name)
	}()

	// pointer
	if t.Kind() == reflect.Ptr {
		str := TypeName(t.Elem())
		if str != "" {
			return "*" + str
		}
		return ""
	}

	// defined types with package path
	if definedName := t.Name(); definedName != "" {
		if pkgPath := t.PkgPath(); pkgPath != "" {
			return pkgPath + "." + definedName
		}
	}

	return ""
}

func Register(t reflect.Type) {
	name := TypeName(t)
	if name == "" {
		panic(fmt.Errorf("not defined type: %v", t))
	}
	registeredNameToType.LoadOrStore(name, t)
	registeredTypeToName.LoadOrStore(t, name)
}
