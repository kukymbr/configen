package gentype

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

func BuildSyntaxMap(pkg *packages.Package) map[string]*ast.StructType {
	m := make(map[string]*ast.StructType)

	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}

			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if st, ok := ts.Type.(*ast.StructType); ok {
					m[ts.Name.Name] = st
				}
			}
		}
	}

	return m
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
