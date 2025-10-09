package gentype

import (
	"context"
	"errors"
)

type Adapter interface {
	Generate(ctx context.Context) (OutputFiles, error)
}

type GenericAdapter struct {
	Source        Source
	OutputOptions OutputOptions
}

func (g *GenericAdapter) Generate(_ context.Context) error {
	return errors.New("not implemented")
}
