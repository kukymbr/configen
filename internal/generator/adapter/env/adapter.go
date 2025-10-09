package env

import (
	"context"
	"strings"

	"github.com/kukymbr/configen/internal/generator/gentype"
)

type Env struct {
	gentype.GenericAdapter

	envs []string
}

func New(sourceStruct gentype.Source, outputOptions gentype.OutputOptions) *Env {
	return &Env{
		GenericAdapter: gentype.GenericAdapter{
			Source:        sourceStruct,
			OutputOptions: outputOptions,
		},

		envs: make([]string, 0),
	}
}

func (g *Env) Generate(ctx context.Context) (gentype.OutputFiles, error) {
	g.collectEnvVars(ctx, g.Source.Struct, "")

	doc := gentype.GetDocComment("#", g.Source.RootStructName, g.Source.RootStructDoc)

	envContent := doc + strings.TrimSpace(strings.Join(g.envs, "\n")) + "\n"

	return gentype.OutputFiles{[]byte(envContent)}, nil
}
