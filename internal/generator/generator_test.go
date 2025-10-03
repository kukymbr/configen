package generator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/kukymbr/configen/internal/formatter"
	"github.com/kukymbr/configen/internal/generator"
)

type generatorGenerateTestCase struct {
	Name                  string
	GetOptFunc            func() generator.Options
	GetContextFunc        func() context.Context
	AssertConstructorFunc func(err error)
	AssertFunc            func(err error)
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
			Name: "with default options",
			GetOptFunc: func() generator.Options {
				return generator.Options{}
			},
			AssertConstructorFunc: func(err error) {
				s.Require().NoError(err)
			},
			AssertFunc: func(err error) {
				s.Require().NoError(err)
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
			Name: "when something unexpectedly happened",
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
	opt := test.GetOptFunc()

	if opt.Formatter == "" {
		opt.Formatter = formatter.Noop
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
		test.AssertFunc(err)
	}
}
