package generator

import (
	"context"
	"errors"
	"fmt"
	"go/types"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/gogetter"
	"github.com/kukymbr/configen/internal/logger"
	"golang.org/x/tools/go/packages"
)

func New(opt Options) (*Generator, error) {
	if err := prepareOptions(&opt); err != nil {
		return nil, err
	}

	logger.Hellof("Hi, this is configen generator.")
	logger.Debugf("Options: " + opt.Debug())

	return &Generator{
		opt: opt,
	}, nil
}

type Generator struct {
	opt Options
}

func (g *Generator) Generate(ctx context.Context) error {
	logger.Debugf("Doing some magic...")

	src, err := g.loadStruct()
	if err != nil {
		return err
	}

	generators := []struct {
		fn  gentype.GeneratorFunc
		out gentype.OutputOptions
	}{
		{fn: generateYAML, out: g.opt.YAML},
		{fn: generateEnv, out: g.opt.Env},
		{fn: gogetter.Generate, out: g.opt.GoGetter},
	}

	for _, gen := range generators {
		if !gen.out.Enable {
			continue
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if err := gen.fn(&src, gen.out); err != nil {
			return err
		}
	}

	logger.Successf("All done.")

	return nil
}

func (g *Generator) loadStruct() (gentype.SourceStruct, error) {
	conf := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedFiles,
		Dir:  g.opt.SourceDir,
	}

	pkgs, err := packages.Load(conf, ".")
	if err != nil {
		return gentype.SourceStruct{}, fmt.Errorf("load package: %w", err)
	}

	if len(pkgs) == 0 {
		return gentype.SourceStruct{}, errors.New("no packages found")
	}

	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(g.opt.StructName)
	if obj == nil {
		return gentype.SourceStruct{}, errors.New("struct not found: " + g.opt.StructName)
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return gentype.SourceStruct{}, fmt.Errorf("%q is not a named type", g.opt.StructName)
	}

	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return gentype.SourceStruct{}, fmt.Errorf("%q is not a struct", g.opt.StructName)
	}

	src := gentype.SourceStruct{
		Package: pkg,
		Struct:  structType,
		Named:   named,

		Name:     g.opt.StructName,
		Doc:      gentype.GetStructDocComment(pkg, g.opt.StructName),
		Comments: gentype.CollectComments(pkg),
	}

	return src, nil
}
