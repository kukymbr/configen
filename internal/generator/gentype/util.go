package gentype

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strings"

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

func ParseNameTag(tagContent string, tagName string, fallback string) string {
	if tagContent == "" {
		return fallback
	}

	st := reflect.StructTag(tagContent)

	nameTag := st.Get(tagName)
	if nameTag == "" {
		return fallback
	}

	parts := strings.Split(nameTag, ",")
	if parts[0] == "-" {
		return ""
	}

	if parts[0] == "" {
		return fallback
	}

	return parts[0]
}

func ParseDefaultValue(tagValue string, tags ...string) string {
	if tagValue == "" {
		return ""
	}

	st := reflect.StructTag(tagValue)

	for _, tag := range tags {
		if v := st.Get(tag); v != "" {
			return v
		}
	}

	return ""
}

func GetUnderlyingStruct(t types.Type) (*types.Struct, *types.Named, bool) {
	switch tt := t.(type) {
	case *types.Pointer:
		return GetUnderlyingStruct(tt.Elem())
	case *types.Named:
		if st, ok := tt.Underlying().(*types.Struct); ok {
			return st, tt, true
		}
	case *types.Struct:
		return tt, nil, true
	}

	return nil, nil, false
}

func DefaultValueForType(t types.Type, value string) string {
	if value != "" {
		return value
	}

	switch tt := t.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return "false"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			return "0"
		case types.String:
			return ""
		default:
			return ""
		}
	}

	return ""
}

func ToPrivateName(name string) string {
	return ToLowerCamel(name)
}

func ToPublicName(name string) string {
	public := ToCamel(name)
	if name != public {
		return public
	}

	return public + "Provider"
}
