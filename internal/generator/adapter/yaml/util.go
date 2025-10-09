package yaml

import (
	"go/types"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"gopkg.in/yaml.v3"
)

func getYAMLBasicNode(t *types.Basic, value string) *yaml.Node {
	var tag string

	val := gentype.DefaultValueForType(t, value)

	switch t.Kind() {
	case types.Bool:
		tag = "!!bool"
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		tag = "!!int"
	case types.Float32, types.Float64:
		tag = "!!float"
	default:
		tag = "!!str"
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: val}
}

func isStructLike(t types.Type) bool {
	switch t.(type) {
	case *types.Struct, *types.Named:
		return true
	}

	return false
}
