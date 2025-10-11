package gogetter

import (
	"bytes"
	"context"
	"go/format"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/logger"
	"github.com/kukymbr/configen/internal/version"
)

type GoGetter struct {
	gentype.GenericAdapter

	collectedStructs map[string]*StructInfo
	collectedImports map[string]struct{}
}

func New(sourceStruct gentype.Source, outputOptions gentype.OutputOptions) *GoGetter {
	return &GoGetter{
		GenericAdapter: gentype.GenericAdapter{
			Source:        sourceStruct,
			OutputOptions: outputOptions,
		},

		collectedStructs: make(map[string]*StructInfo),
		collectedImports: make(map[string]struct{}),
	}
}

func (g *GoGetter) Generate(ctx context.Context) (gentype.OutputFiles, error) {
	if g.OutputOptions.TargetPackageName == "" {
		g.OutputOptions.TargetPackageName = packageNameFromID(g.Source.Package.ID)
	}

	g.processStruct(ctx, g.Source.Named, g.Source.Struct, g.OutputOptions.TargetStructName, false)

	tplData := tplData{
		Structs:          g.collectedStructs,
		Imports:          g.getImports(),
		PackageName:      g.OutputOptions.TargetPackageName,
		Version:          version.GetVersion(),
		TargetStructName: g.OutputOptions.TargetStructName,
		SourceStructName: g.Source.RootStructName,
	}

	var buf bytes.Buffer
	if err := executeTemplate(&buf, tplData); err != nil {
		return nil, err
	}

	content := buf.Bytes()

	formatted, err := format.Source(content)
	if err == nil {
		content = formatted
	} else {
		logger.Warningf("Failed to format generated code: %s", err.Error())
	}

	return gentype.OutputFiles{content}, nil
}
