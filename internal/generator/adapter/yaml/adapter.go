package yaml

import (
	"context"
	"fmt"

	"github.com/kukymbr/configen/internal/generator/gentype"
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

func (g *YAML) Generate(ctx context.Context) (gentype.OutputFiles, error) {
	yamlNode := g.structToYAMLNode(ctx, g.Source.Struct)

	data, err := yaml.Marshal(yamlNode)
	if err != nil {
		return nil, fmt.Errorf("marshal YAML nodes: %w", err)
	}

	doc := gentype.GetDocComment("#", g.Source.RootStructName, g.Source.RootStructDoc)

	data = append([]byte(doc), data...)

	return gentype.OutputFiles{data}, nil
}
