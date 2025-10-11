package gogetter

import (
	"context"
	"fmt"
	"go/token"
	"go/types"
	"slices"

	"github.com/kukymbr/configen/internal/generator/gentype"
)

func (g *GoGetter) processStruct(
	ctx context.Context,
	named *types.Named,
	st *types.Struct,
	targetStructName string,
	isAnon bool,
) *StructInfo {
	ctx = gentype.ContextIncRecursionDepth(ctx)
	gentype.ContextMustValidateRecursionDepth(ctx, "Go generator (processStruct)")

	if ctx.Err() != nil {
		return nil
	}

	syntaxMap := g.Source.SyntaxMap

	if targetStructName == "" {
		targetStructName = gentype.ToPublicName(named.Obj().Name())
	}

	if info, exists := g.collectedStructs[targetStructName]; exists {
		return info
	}

	info := &StructInfo{
		Name:             targetStructName,
		SourceStructName: named.Obj().Name(),
		IsAnonymous:      isAnon,
	}

	if syn, ok := syntaxMap[named.Obj().Name()]; ok {
		info.Doc = docComment(syn)
	}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		ft := field.Type()

		if fieldInfo := g.processField(ctx, field, ft, named.Obj().Name(), targetStructName, i); fieldInfo != nil {
			info.Fields = append(info.Fields, fieldInfo...)
		}
	}

	g.collectedStructs[targetStructName] = info

	return info
}

//nolint:cyclop
func (g *GoGetter) processField(
	ctx context.Context,
	field *types.Var,
	ft types.Type,
	sourceStructName string,
	targetStructName string,
	fieldIndex int,
) []FieldInfo {
	// Handle anonymous embedded structs by flattening their fields
	if field.Anonymous() {
		return g.processAnonymousField(ctx, field, targetStructName)
	}

	if !field.Exported() {
		return nil
	}

	var structInfo *StructInfo

	typeName := g.formatTypeName(ft)
	processed := false

	if nt, ok := ft.(*types.Named); ok {
		if _, ok := nt.Underlying().(*types.Struct); ok && g.isTargetPackage(nt.Obj().Pkg()) {
			typeName = gentype.ToPublicName(nt.Obj().Name())
			structInfo = g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), typeName, false)
			processed = true
		}
	}

	if pt, ok := ft.(*types.Pointer); ok && !processed {
		if nt, ok := pt.Elem().(*types.Named); ok {
			if _, ok := nt.Underlying().(*types.Struct); ok && g.isTargetPackage(nt.Obj().Pkg()) {
				pubName := gentype.ToPublicName(nt.Obj().Name())
				typeName = "*" + pubName

				structInfo = g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), pubName, false)
				processed = true
			}
		}
	}

	if stt, ok := ft.(*types.Struct); ok && !processed {
		structInfo = g.processStruct(
			ctx,
			g.anonStructToNamed(stt, targetStructName, field),
			stt, "",
			true,
		)
	}

	return []FieldInfo{{
		Name:       gentype.ToPrivateName(field.Name()),
		ExportName: field.Name(),
		TypeName:   typeName,
		Comment:    g.Source.GetStructFieldComment(sourceStructName, fieldIndex),
		IsStruct:   structInfo != nil,
		StructInfo: structInfo,
	}}
}

func (g *GoGetter) processAnonymousField(
	ctx context.Context,
	field *types.Var,
	targetStructName string,
) []FieldInfo {
	ft := field.Type()

	if nt, ok := ft.(*types.Named); ok {
		if _, ok := nt.Underlying().(*types.Struct); ok {
			g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), "", true)

			embedded := g.collectedStructs[gentype.ToPublicName(nt.Obj().Name())]

			return embedded.Fields
		}
	}

	if stt, ok := ft.(*types.Struct); ok {
		named := g.anonStructToNamed(stt, targetStructName, field)

		g.processStruct(
			ctx,
			g.anonStructToNamed(stt, targetStructName, field),
			stt,
			"",
			true,
		)

		embedded := g.collectedStructs[gentype.ToPublicName(named.Obj().Name())]

		return embedded.Fields
	}

	return nil
}

func (g *GoGetter) anonStructToNamed(st *types.Struct, targetStructName string, field *types.Var) *types.Named {
	anonName := fmt.Sprintf("%s%s", targetStructName, gentype.ToCamel(field.Name()))

	return types.NewNamed(types.NewTypeName(
		token.NoPos, nil, anonName, nil,
	), st, nil)
}

func (g *GoGetter) getImports() []string {
	imports := make([]string, 0, len(g.collectedImports))

	for imp := range g.collectedImports {
		imports = append(imports, imp)
	}

	slices.Sort(imports)

	return imports
}

//nolint:cyclop
func (g *GoGetter) formatTypeName(t types.Type) string {
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name()
	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()

		if pkg != nil && pkg.Name() != "main" && !g.isTargetPackage(pkg) {
			g.registerImport(pkg.Path())

			return pkg.Name() + "." + obj.Name()
		}

		return obj.Name()
	case *types.Slice:
		return "[]" + g.formatTypeName(tt.Elem())
	case *types.Pointer:
		return "*" + g.formatTypeName(tt.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", tt.Len(), g.formatTypeName(tt.Elem()))
	case *types.Map:
		return fmt.Sprintf(
			"map[%s]%s",
			g.formatTypeName(tt.Key()),
			g.formatTypeName(tt.Elem()),
		)
	default:
		return "any"
	}
}

func (g *GoGetter) registerImport(imp string) {
	if g.isTargetPackage(imp) {
		return
	}

	g.collectedImports[imp] = struct{}{}
}

func (g *GoGetter) isTargetPackage(pkg any) bool {
	pkgName := gentype.ParsePackageName(pkg)

	return pkgName == g.OutputOptions.TargetPackageName
}
