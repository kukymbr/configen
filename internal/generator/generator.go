package generator

import (
	"context"
	"errors"
	"fmt"
	"go/types"

	"github.com/kukymbr/configen/internal/generator/adapter/gogetter"
	"github.com/kukymbr/configen/internal/generator/adapter/yaml"
	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/utils"
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
		adapter func(out gentype.OutputOptions) gentype.Adapter
		out     gentype.OutputOptions
	}{
		{
			adapter: func(out gentype.OutputOptions) gentype.Adapter {
				return yaml.New(src, out)
			},
			out: g.opt.YAML,
		},
		// TODO
		//{adapter: generateEnv, out: g.opt.Env},
		{
			adapter: func(out gentype.OutputOptions) gentype.Adapter {
				return gogetter.New(src, out)
			},
			out: g.opt.GoGetter,
		},
	}

	for _, gen := range generators {
		if !gen.out.Enable {
			continue
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		adapter := gen.adapter(gen.out)

		// TODO: run in routines
		files, err := adapter.Generate(ctx)
		if err != nil {
			return err
		}

		for _, content := range files {
			if err := utils.WriteFile(content, gen.out.Path); err != nil {
				return err
			}
		}
	}

	logger.Successf("All done.")

	return nil
}

func (g *Generator) loadStruct() (gentype.Source, error) {
	conf := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedFiles,
		Dir:  g.opt.SourceDir,
	}

	pkgs, err := packages.Load(conf, ".")
	if err != nil {
		return gentype.Source{}, fmt.Errorf("load package: %w", err)
	}

	if len(pkgs) == 0 {
		return gentype.Source{}, errors.New("no packages found")
	}

	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(g.opt.StructName)
	if obj == nil {
		return gentype.Source{}, errors.New("struct not found: " + g.opt.StructName)
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return gentype.Source{}, fmt.Errorf("%q is not a named type", g.opt.StructName)
	}

	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return gentype.Source{}, fmt.Errorf("%q is not a struct", g.opt.StructName)
	}

	return gentype.NewSource(pkg, g.opt.StructName, named, structType), nil
}
