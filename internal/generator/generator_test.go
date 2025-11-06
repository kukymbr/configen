package generator_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kukymbr/configen/internal/generator"
	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/stretchr/testify/suite"
)

const (
	// Using the example config file make changes once in example file
	// and validate it too into the tests.
	givenSourceDir  = "../../example"
	givenStructName = "config"
)

type generatorGenerateTestCase struct {
	Name                  string
	GetOptFunc            func() generator.Options
	GetContextFunc        func() context.Context
	AssertConstructorFunc func(err error)
	AssertFunc            func(opt generator.Options, err error)
}

func TestGenerator(t *testing.T) {
	suite.Run(t, &GeneratorSuite{})
}

type GeneratorSuite struct {
	suite.Suite
}

func (s *GeneratorSuite) SetupSuite() {}

func (s *GeneratorSuite) TearDownSuite() {}

func (s *GeneratorSuite) TestGenerator_PositiveCases() {
	tests := []generatorGenerateTestCase{
		{
			Name: "generate all",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					YAML: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
					Env: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
					GoGetter: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(opt generator.Options, err error) {
				s.Require().NoError(err)

				s.assertContent(opt.YAML.Path, "config.yaml")
				s.assertContent(opt.Env.Path, "config.env")
				s.assertContent(opt.GoGetter.Path, "config.gen.go")
			},
		},
		{
			Name: "generate local",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					YAML: gentype.OutputOptions{
						Enable:          true,
						Path:            s.getTargetPath(),
						Tag:             "local",
						DefaultValueTag: "localDefault",
					},
					Env: gentype.OutputOptions{
						Enable:          true,
						Path:            s.getTargetPath(),
						DefaultValueTag: "localDefault",
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(opt generator.Options, err error) {
				s.Require().NoError(err)

				s.assertContent(opt.YAML.Path, "local.yaml")
				s.assertContent(opt.Env.Path, "local.env")
			},
		},
		{
			Name: "no generator enabled",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					YAML: gentype.OutputOptions{
						Enable: false,
						Path:   s.getTargetPath(),
					},
					Env: gentype.OutputOptions{
						Enable: false,
						Path:   s.getTargetPath(),
					},
					GoGetter: gentype.OutputOptions{
						Enable: false,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(opt generator.Options, err error) {
				s.Require().NoError(err)

				s.NoFileExists(opt.YAML.Path)
				s.NoFileExists(opt.Env.Path)
				s.NoFileExists(opt.GoGetter.Path)
			},
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.runGeneratorGenerateTest(test)
		})
	}
}

func (s *GeneratorSuite) TestGenerator_NegativeCases() {
	tests := []generatorGenerateTestCase{
		{
			Name: "empty options given",
			GetOptFunc: func() generator.Options {
				return generator.Options{}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "unknown struct given",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: "UnknownStruct",
					Env: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(opt generator.Options, err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "invalid struct format",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: "***",
					GoGetter: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "empty struct name",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: " ",
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "invalid source dir",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					SourceDir:  "testdata/unknown",
					GoGetter: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "invalid source not a dir",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					SourceDir:  "testdata/config.go",
					GoGetter: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
			},
		},
		{
			Name: "context is canceled",
			GetOptFunc: func() generator.Options {
				return generator.Options{
					StructName: givenStructName,
					YAML: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
					Env: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
					GoGetter: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath(),
					},
				}
			},
			GetContextFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return ctx
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(opt generator.Options, err error) {
				s.Require().ErrorIs(err, context.Canceled)

				s.NoFileExists(opt.YAML.Path)
				s.NoFileExists(opt.Env.Path)
				s.NoFileExists(opt.GoGetter.Path)
			},
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.runGeneratorGenerateTest(test)
		})
	}
}

func (s *GeneratorSuite) runGeneratorGenerateTest(test generatorGenerateTestCase) {
	s.T().Helper()

	opt := test.GetOptFunc()

	if opt.SourceDir == "" {
		opt.SourceDir = givenSourceDir
	}

	gen, err := generator.New(opt)
	test.AssertConstructorFunc(err)

	if err != nil {
		return
	}

	ctx := s.T().Context()
	if test.GetContextFunc != nil {
		ctx = test.GetContextFunc()
	}

	err = gen.Generate(ctx)
	if test.AssertFunc != nil {
		test.AssertFunc(opt, err)
	}
}

func (s *GeneratorSuite) assertContent(filename string, expectedName string) {
	s.T().Helper()

	actual, err := os.ReadFile(filename)
	s.Require().NoError(err)

	expectedPath := filepath.Join("testdata/expected", expectedName)

	expected, err := os.ReadFile(expectedPath)
	s.Require().NoError(err)

	s.Require().Equal(string(expected), string(actual))
}

func (s *GeneratorSuite) getTargetPath() string {
	s.T().Helper()

	name := fmt.Sprintf(
		"%s_%d-%d.tmp",
		strings.ReplaceAll(s.T().Name(), "/", "."),
		time.Now().UnixNano(),
		rand.Uint(),
	)
	path := filepath.Join("testdata/target", name)

	s.T().Cleanup(func() {
		_ = os.RemoveAll(path)
	})

	return path
}
