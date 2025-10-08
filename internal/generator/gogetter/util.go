package gogetter

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
)

func packageNameFromID(id string) string {
	parts := strings.Split(id, "/")

	return parts[len(parts)-1]
}

func docComment(st *ast.StructType) string {
	if st.Fields == nil {
		return ""
	}

	for _, f := range st.Fields.List {
		if f.Doc != nil {
			return f.Doc.Text()
		}
	}

	return ""
}

func fieldComment(field *ast.Field) string {
	if field.Doc != nil {
		return field.Doc.Text()
	}

	if field.Comment != nil {
		return field.Comment.Text()
	}

	return ""
}

func toPrivateName(s string) string {
	return strcase.ToLowerCamel(s)
}

func toPublicName(name string) string {
	return strcase.ToCamel(name)
}

//nolint:cyclop
func formatTypeName(t types.Type, targetPackageName string, imports *[]string) string {
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name()
	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()
		pkgName := packageNameFromID(pkg.Path())

		if pkg != nil && pkg.Name() != "main" && pkgName != targetPackageName {
			*imports = append(*imports, pkg.Path())

			return pkg.Name() + "." + obj.Name()
		}

		return obj.Name()
	case *types.Slice:
		return "[]" + formatTypeName(tt.Elem(), targetPackageName, imports)
	case *types.Pointer:
		return "*" + formatTypeName(tt.Elem(), targetPackageName, imports)
	case *types.Array:
		return fmt.Sprintf("[%d]%s", tt.Len(), formatTypeName(tt.Elem(), targetPackageName, imports))
	case *types.Map:
		return fmt.Sprintf(
			"map[%s]%s",
			formatTypeName(tt.Key(), targetPackageName, imports),
			formatTypeName(tt.Elem(), targetPackageName, imports),
		)
	default:
		return "any"
	}
}
