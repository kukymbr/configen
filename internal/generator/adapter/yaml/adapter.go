package yaml

import (
	"context"
	"fmt"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/logger"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	gentype.GenericAdapter

	visited map[string]struct{}
}

func New(sourceStruct gentype.Source, outputOptions gentype.OutputOptions) *YAML {
	return &YAML{
		GenericAdapter: gentype.GenericAdapter{
			Source:        sourceStruct,
			OutputOptions: outputOptions,
		},

		visited: make(map[string]struct{}),
	}
}

func (g *YAML) Name() string {
	return "YAMLGenerator"
}

func (g *YAML) Generate(ctx context.Context) (gentype.OutputFiles, error) {
	yamlNode := g.structToYAMLNode(ctx, g.Source.Struct)

	g.debugf("marshaling...")

	data, err := yaml.Marshal(yamlNode)
	if err != nil {
		return nil, fmt.Errorf("marshal YAML nodes: %w", err)
	}

	g.debugf("writing head comment...")

	doc := gentype.GetDocComment("#", g.Source.RootStructName, g.Source.RootStructDoc)

	data = append([]byte(doc), data...)

	return gentype.OutputFiles{data}, nil
}

func (g *YAML) debugf(format string, args ...any) {
	logger.Debugf(g.Name()+": "+format, args...)
}
