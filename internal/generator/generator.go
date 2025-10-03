package generator

import (
	"context"

	"github.com/kukymbr/configen/internal/formatter"
	"github.com/kukymbr/configen/internal/logger"
)

func New(opt Options) (*Generator, error) {
	if err := prepareOptions(&opt); err != nil {
		return nil, err
	}

	f, err := formatter.Factory(opt.Formatter)
	if err != nil {
		return nil, err
	}

	logger.Hellof("Hi, this is configen generator.")
	logger.Debugf("Options: " + opt.Debug())

	return &Generator{
		opt:       opt,
		formatter: f,
	}, nil
}

type Generator struct {
	opt       Options
	formatter formatter.Formatter
}

func (g *Generator) Generate(ctx context.Context) error {
	logger.Debugf("Doing some magic...")

	// TODO: do some magic.

	logger.Successf("All done.")

	return nil
}

func (g *Generator) format(ctx context.Context, content []byte) []byte {
	logger.Debugf("Formatting generated code...")

	formatted, err := g.formatter.Format(ctx, content)
	if err != nil {
		logger.Warningf("Failed to format generated code: %s", err.Error())

		return content
	}

	return formatted
}
