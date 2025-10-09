package generator

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/utils"
	"github.com/kukymbr/configen/internal/logger"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
)

func generateYAML(src *gentype.Source, out gentype.OutputOptions) error {
	yamlNode := structToYAMLNode(src.Package, src.Struct, src.CommentsMap, nil, out.Tag)

	data, err := yaml.Marshal(yamlNode)
	if err != nil {
		return fmt.Errorf("marshal YAML nodes: %w", err)
	}

	doc := gentype.GetDocComment("#", src.RootStructName, src.RootStructDoc)

	data = append([]byte(doc), data...)

	if err := utils.WriteFile(data, out.Path); err != nil {
		return err
	}

	logger.Successf("Generated %s file", out.Path)

	return nil
}

func structToYAMLNode(
	pkg *packages.Package,
	st *types.Struct,
	comments map[token.Pos]string,
	visited map[string]bool,
	tagName string,
) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Exported() {
			continue
		}

		tag := st.Tag(i)

		yamlName := parseNameTag(tag, tagName, field.Name())
		if yamlName == "" {
			continue
		}

		value := parseDefaultValue(tag, valueTagsYAML...)
		comment := comments[field.Pos()]
		ft := field.Type()

		if field.Anonymous() {
			if stt, ok := underlyingStruct(ft); ok {
				embedded := structToYAMLNode(pkg, stt, comments, visited, tagName)
				node.Content = append(node.Content, embedded.Content...)

				continue
			}
		}

		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: yamlName}
		valNode := typeToYAMLNode(pkg, ft, comments, visited, value, tagName)

		if comment != "" {
			keyNode.HeadComment = comment
		}

		node.Content = append(node.Content, keyNode, valNode)
	}

	return node
}

func getYAMLBasicNode(t *types.Basic, value string) *yaml.Node {
	var tag string

	val := defaultValueForType(t, value)

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

//nolint:cyclop,funlen,gocognit
func typeToYAMLNode(
	pkg *packages.Package,
	t types.Type,
	comments map[token.Pos]string,
	visited map[string]bool,
	value string,
	tagName string,
) *yaml.Node {
	if visited == nil {
		visited = make(map[string]bool)
	}

	if isTextMarshaler(t) || isStringer(t) {
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value,
			Tag:   "!!str",
		}
	}

	switch tt := t.(type) {
	case *types.Basic:
		return getYAMLBasicNode(tt, value)
	case *types.Pointer:
		return typeToYAMLNode(pkg, tt.Elem(), comments, visited, value, tagName)
	case *types.Slice, *types.Array:
		elemType := tt.(interface{ Elem() types.Type }).Elem()

		seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

		if value != "" && !isStructLike(elemType) {
			for _, v := range strings.Split(value, ",") {
				elemNode := typeToYAMLNode(pkg, elemType, comments, visited, v, tagName)
				seq.Content = append(seq.Content, elemNode)
			}
		} else {
			seq.Content = append(seq.Content, typeToYAMLNode(pkg, elemType, comments, visited, "", tagName))
		}

		return seq
	case *types.Map:
		if value != "" {
			m := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

			for _, p := range strings.Split(value, ",") {
				kv := strings.SplitN(p, "=", 2)
				k := kv[0]
				v := ""

				if len(kv) == 2 {
					v = kv[1]
				}

				m.Content = append(m.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: k},
					&yaml.Node{Kind: yaml.ScalarNode, Value: v})
			}

			return m
		}

		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	case *types.Named:
		if st, ok := tt.Underlying().(*types.Struct); ok {
			key := tt.String()
			if visited[key] {
				return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map",
					Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "<recursive>"}}}
			}

			visited[key] = true

			return structToYAMLNode(pkg, st, comments, visited, tagName)
		}

		return typeToYAMLNode(pkg, tt.Underlying(), comments, visited, value, tagName)
	case *types.Struct:
		key := tt.String()
		if visited[key] {
			return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map",
				Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "<recursive>"}}}
		}

		visited[key] = true

		return structToYAMLNode(pkg, tt, comments, visited, tagName)
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: ""}
}
