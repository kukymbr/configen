package gogetter

import (
	"go/ast"
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

func toPrivateName(s string) string {
	return strcase.ToLowerCamel(s)
}

func toPublicName(name string) string {
	return strcase.ToCamel(name)
}
