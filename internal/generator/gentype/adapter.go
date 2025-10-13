package gentype

import (
	"context"
	"errors"
)

type Adapter interface {
	Name() string
	Generate(ctx context.Context) (OutputFiles, error)
}

type GenericAdapter struct {
	Source        Source
	OutputOptions OutputOptions
}

func (g *GenericAdapter) Name() string {
	return "GenericAdapter"
}

func (g *GenericAdapter) Generate(_ context.Context) (OutputFiles, error) {
	return nil, errors.New("not implemented")
}
