package sb

import "reflect"

func TypeName(t reflect.Type) string {

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
