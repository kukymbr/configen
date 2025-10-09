package utils

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func ParsePackageName(pkg any) string {
	switch p := pkg.(type) {
	case *types.PkgName:
		return packageNameFromPath(p.Name())
	case *types.Package:
		for _, opt := range [2]string{p.Name(), p.Path()} {
			if name := packageNameFromPath(opt); name != "" {
				return name
			}
		}
	case *packages.Package:
		for _, opt := range [2]string{p.Name, p.ID} {
			if name := packageNameFromPath(opt); name != "" {
				return name
			}
		}
	case *types.Named:
		return packageNameFromPath(p.Obj().Pkg().Path())
	case string:
		if name := packageNameFromType(p); name != "" {
			return name
		}

		return packageNameFromPath(p)
	}

	return ""
}

func packageNameFromPath(id string) string {
	parts := strings.Split(id, "/")

	return parts[len(parts)-1]
}

func packageNameFromType(t string) string {
	parts := strings.Split(t, ".")
	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}
