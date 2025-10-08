package gogetter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/logger"
	"github.com/kukymbr/configen/internal/utils"
	"github.com/kukymbr/configen/internal/version"
	"golang.org/x/tools/go/packages"
)

func Generate(src *gentype.SourceStruct, out gentype.OutputOptions) error {
	if out.TargetPackageName == "" {
		out.TargetPackageName = packageNameFromID(src.Package.ID)
	}

	collected := map[string]*StructInfo{}
	syntaxMap := buildSyntaxMap(src.Package)
	imports := make([]string, 0)
	targetStructName := toPublicName(out.TargetStructName)

	processStruct(src.Named, src.Struct, syntaxMap, collected, &imports, targetStructName, out.TargetPackageName)

	tplData := tplData{
		Structs:          collected,
		Imports:          filterImports(out.TargetPackageName, imports),
		PackageName:      out.TargetPackageName,
		Version:          version.GetVersion(),
		TargetStructName: targetStructName,
		SourceStructName: src.Name,
	}

	var buf bytes.Buffer
	if err := executeTemplate(&buf, tplData); err != nil {
		return err
	}

	content := buf.Bytes()

	formatted, err := format.Source(content)
	if err == nil {
		content = formatted
	}

	if err := utils.WriteFile(content, out.Path); err != nil {
		return err
	}

	logger.Successf("Generated %s file", out.Path)

	return nil
}

func buildSyntaxMap(pkg *packages.Package) map[string]*ast.StructType {
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

//nolint:gocognit,cyclop,funlen
func processStruct(
	named *types.Named,
	st *types.Struct,
	syntaxMap map[string]*ast.StructType,
	collected map[string]*StructInfo,
	imports *[]string,
	targetStructName string,
	targetPackageName string,
) {
	if targetStructName == "" {
		targetStructName = toPublicName(named.Obj().Name())
	}

	if _, exists := collected[targetStructName]; exists {
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

		// Handle anonymous embedded structs by flattening their fields
		if field.Anonymous() {
			if nt, ok := ft.(*types.Named); ok {
				if _, ok := nt.Underlying().(*types.Struct); ok {
					processStruct(
						nt, nt.Underlying().(*types.Struct),
						syntaxMap, collected, imports,
						"", targetPackageName,
					)
					embedded := collected[toPublicName(nt.Obj().Name())]

					info.Fields = append(info.Fields, embedded.Fields...)

					continue
				}
			}

			if stt, ok := ft.(*types.Struct); ok {
				anonName := fmt.Sprintf("%s_anon_%d", targetStructName, i)

				processStruct(
					types.NewNamed(types.NewTypeName(0, nil, anonName, nil), stt, nil),
					stt,
					syntaxMap,
					collected,
					imports,
					"",
					targetPackageName,
				)

				embedded := collected[toPublicName(anonName)]

				info.Fields = append(info.Fields, embedded.Fields...)

				continue
			}
		}

		if !field.Exported() {
			continue
		}

		typeName := formatTypeName(ft, targetPackageName, imports)
		isStruct := false

		if nt, ok := ft.(*types.Named); ok {
			if _, ok := nt.Underlying().(*types.Struct); ok {
				isStruct = true
				typeName = toPublicName(nt.Obj().Name())

				processStruct(
					nt, nt.Underlying().(*types.Struct),
					syntaxMap, collected, imports,
					"", targetPackageName,
				)
			}
		}

		if pt, ok := ft.(*types.Pointer); ok {
			if nt, ok := pt.Elem().(*types.Named); ok {
				if _, ok := nt.Underlying().(*types.Struct); ok {
					isStruct = true
					typeName = "*" + toPublicName(nt.Obj().Name())

					processStruct(
						nt, nt.Underlying().(*types.Struct),
						syntaxMap, collected, imports,
						"", targetPackageName,
					)
				}
			}
		}

		comment := ""

		if syn, ok := syntaxMap[named.Obj().Name()]; ok {
			if i < len(syn.Fields.List) {
				comment = fieldComment(syn.Fields.List[i])
			}
		}

		fieldInfo := FieldInfo{
			Name:       toPrivateName(field.Name()),
			ExportName: field.Name(),
			TypeName:   typeName,
			Comment:    comment,
			IsStruct:   isStruct,
		}

		info.Fields = append(info.Fields, fieldInfo)
	}

	collected[targetStructName] = info
}

func filterImports(targetPackageName string, imports []string) []string {
	filtered := make([]string, 0, len(imports))

	for _, imp := range imports {
		pkg := packageNameFromID(imp)
		if pkg == targetPackageName {
			continue
		}

		filtered = append(filtered, imp)
	}

	return filtered
}
