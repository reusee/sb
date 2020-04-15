// +build ignore

package main

import (
	"fmt"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"
)

var pt = fmt.Printf

func main() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypesInfo |
			packages.NeedTypes |
			packages.NeedName,
	}, "std")
	if err != nil {
		panic(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return
	}

	stat := make(map[string]int)
	var names []string
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			t := obj.Type().Underlying()
			basic, ok := t.(*types.Basic)
			if !ok {
				continue
			}
			name := basic.String()
			if _, ok := stat[name]; !ok {
				names = append(names, name)
			}
			stat[name]++
		}
	})

	sort.Slice(names, func(i, j int) bool {
		return stat[names[i]] > stat[names[j]]
	})
	for _, n := range names {
		pt("%s %d\n", n, stat[n])
	}

}
