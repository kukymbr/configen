package generator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kukymbr/configen/internal/generator"
	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/stretchr/testify/suite"
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
					StructName: "Config",
					YAML: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath("test1.yaml"),
					},
					Env: gentype.OutputOptions{
						Enable: true,
						Path:   s.getTargetPath("test1.env"),
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
			Name: "required options missing",
			GetOptFunc: func() generator.Options {
				return generator.Options{}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().Error(err)
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
	opt.SourceDir = "testdata"

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

func (s *GeneratorSuite) getTargetPath(name string) string {
	s.T().Helper()

	path := filepath.Join("testdata/target", name)

	s.T().Cleanup(func() {
		_ = os.RemoveAll(path)
	})

	return path
}
