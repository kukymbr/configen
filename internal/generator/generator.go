package generator

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"go/types"

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

func (g *Generator) Generate(_ context.Context) error {
	logger.Debugf("Doing some magic...")

	src, err := g.loadStruct()
	if err != nil {
		return err
	}

	generators := []struct {
		fn  generatorFunc
		out OutputOptions
	}{
		{fn: generateYAML, out: g.opt.YAML},
		{fn: generateEnv, out: g.opt.Env},
	}

	for _, gen := range generators {
		if !gen.out.Enable {
			continue
		}

		if err := gen.fn(&src, gen.out.Path); err != nil {
			return err
		}
	}

	logger.Successf("All done.")

	return nil
}

func (g *Generator) loadStruct() (sourceStruct, error) {
	conf := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedFiles,
		Dir:  g.opt.SourceDir,
	}

	pkgs, err := packages.Load(conf, ".")
	if err != nil {
		return sourceStruct{}, fmt.Errorf("load package: %w", err)
	}

	if len(pkgs) == 0 {
		return sourceStruct{}, errors.New("no packages found")
	}

	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(g.opt.StructName)
	if obj == nil {
		return sourceStruct{}, errors.New("struct not found: " + g.opt.StructName)
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return sourceStruct{}, fmt.Errorf("%q is not a named type", g.opt.StructName)
	}

	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return sourceStruct{}, fmt.Errorf("%q is not a struct", g.opt.StructName)
	}

	src := sourceStruct{
		pkg:      pkg,
		st:       structType,
		name:     g.opt.StructName,
		doc:      getStructDocComment(pkg, g.opt.StructName),
		comments: collectComments(pkg),
	}

	return src, nil
}

type sourceStruct struct {
	pkg      *packages.Package
	st       *types.Struct
	name     string
	doc      string
	comments map[token.Pos]string
}

type generatorFunc func(src *sourceStruct, target string) error
