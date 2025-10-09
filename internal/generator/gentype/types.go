package gentype

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Source struct {
	Package *packages.Package
	Struct  *types.Struct
	Named   *types.Named

	RootStructName string
	RootStructDoc  string

	CommentsMap map[token.Pos]string
	SyntaxMap   map[string]*ast.StructType
}

func NewSource(pkg *packages.Package, structName string, named *types.Named, st *types.Struct) Source {
	return Source{
		Package:        pkg,
		Struct:         st,
		Named:          named,
		RootStructName: structName,
		RootStructDoc:  GetStructDocComment(pkg, structName),
		CommentsMap:    BuildCommentsMap(pkg),
		SyntaxMap:      BuildSyntaxMap(pkg),
	}
}

func (s *Source) GetStructFieldComment(structName string, fieldIndex int) string {
	if syn, ok := s.SyntaxMap[structName]; ok && fieldIndex < len(syn.Fields.List) {
		return fieldComment(syn.Fields.List[fieldIndex])
	}

	return ""
}

type OutputOptions struct {
	// Enable is a flag to enable an output.
	Enable bool

	// Path is a target file path.
	Path string

	// Tag is a field names tag.
	Tag string

	TargetStructName string

	TargetPackageName string
}

type OutputFiles [][]byte

type Nullable[T any] struct {
	value *T
}

func NewNullable[T any](value ...T) Nullable[T] {
	v := Nullable[T]{}

	if len(value) > 0 {
		v.Set(value[0])
	}

	return v
}

func (v *Nullable[T]) IsSet() bool {
	return v.value != nil
}

func (v *Nullable[T]) Value() T {
	var empty T
	if v.value == nil {
		return empty
	}

	return *v.value
}

func (v *Nullable[T]) Set(value T) {
	v.value = &value
}

func (v *Nullable[T]) Unset() {
	v.value = nil
}
