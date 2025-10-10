package gogetter

import (
	"context"
	"fmt"
	"go/types"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/utils"
)

func (g *GoGetter) processStruct(ctx context.Context, named *types.Named, st *types.Struct, targetStructName string) {
	ctx = gentype.ContextIncRecursionDepth(ctx)
	gentype.ContextMustValidateRecursionDepth(ctx, "Go generator (processStruct)")

	if ctx.Err() != nil {
		return
	}

	syntaxMap := g.Source.SyntaxMap

	if targetStructName == "" {
		targetStructName = toPublicName(named.Obj().Name())
	}

	if _, exists := g.collectedStructs[targetStructName]; exists {
		return
	}

	info := &StructInfo{
		Name:             targetStructName,
		SourceStructName: named.Obj().Name(),
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
}

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
		return g.processAnonymousField(ctx, field.Type(), targetStructName, fieldIndex)
	}

	if !field.Exported() {
		return nil
	}

	typeName := g.formatTypeName(ft)
	isStruct := false

	if nt, ok := ft.(*types.Named); ok {
		if _, ok := nt.Underlying().(*types.Struct); ok && g.isTargetPackage(nt.Obj().Pkg()) {
			isStruct = true
			typeName = toPublicName(nt.Obj().Name())

			g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), typeName)
		}
	}

	if pt, ok := ft.(*types.Pointer); ok {
		if nt, ok := pt.Elem().(*types.Named); ok {
			if _, ok := nt.Underlying().(*types.Struct); ok && g.isTargetPackage(nt.Obj().Pkg()) {
				isStruct = true
				pubName := toPublicName(nt.Obj().Name())
				typeName = "*" + pubName

				g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), pubName)
			}
		}
	}

	return []FieldInfo{{
		Name:       toPrivateName(field.Name()),
		ExportName: field.Name(),
		TypeName:   typeName,
		Comment:    g.Source.GetStructFieldComment(sourceStructName, fieldIndex),
		IsStruct:   isStruct,
	}}
}

func (g *GoGetter) processAnonymousField(
	ctx context.Context,
	ft types.Type,
	targetStructName string,
	fieldIndex int,
) []FieldInfo {
	if nt, ok := ft.(*types.Named); ok {
		if _, ok := nt.Underlying().(*types.Struct); ok {
			g.processStruct(ctx, nt, nt.Underlying().(*types.Struct), "")

			embedded := g.collectedStructs[toPublicName(nt.Obj().Name())]

			return embedded.Fields
		}
	}

	if stt, ok := ft.(*types.Struct); ok {
		anonName := fmt.Sprintf("%s_anon_%d", targetStructName, fieldIndex)

		g.processStruct(
			ctx,
			types.NewNamed(types.NewTypeName(0, nil, anonName, nil), stt, nil),
			stt,
			"",
		)

		embedded := g.collectedStructs[toPublicName(anonName)]

		return embedded.Fields
	}

	return nil
}

func (g *GoGetter) getImports() []string {
	imports := make([]string, 0, len(g.collectedImports))

	for imp := range g.collectedImports {
		imports = append(imports, imp)
	}

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
	pkgName := utils.ParsePackageName(pkg)

	return pkgName == g.OutputOptions.TargetPackageName
}
