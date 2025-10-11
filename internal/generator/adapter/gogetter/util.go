package gogetter

import (
	"go/ast"
	"strings"
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
