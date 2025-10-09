package yaml

import (
	"context"
	"fmt"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/utils"
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

func (g *YAML) Generate(ctx context.Context) error {
	yamlNode := g.structToYAMLNode(ctx, g.Source.Struct)

	data, err := yaml.Marshal(yamlNode)
	if err != nil {
		return fmt.Errorf("marshal YAML nodes: %w", err)
	}

	doc := gentype.GetDocComment("#", g.Source.RootStructName, g.Source.RootStructDoc)

	data = append([]byte(doc), data...)

	if err := utils.WriteFile(data, g.OutputOptions.Path); err != nil {
		return err
	}

	logger.Successf("Generated %s file", g.OutputOptions.Path)

	return nil
}
