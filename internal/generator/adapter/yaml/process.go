package yaml

import (
	"context"
	"go/types"
	"strings"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"gopkg.in/yaml.v3"
)

func (g *YAML) structToYAMLNode(ctx context.Context, st *types.Struct) *yaml.Node {
	g.debugf("converting struct to yaml node")

	ctx = gentype.ContextIncRecursionDepth(ctx)
	gentype.ContextMustValidateRecursionDepth(ctx, "YAML generator (structToYAMLNode)")

	if ctx.Err() != nil {
		return nil
	}

	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		tag := st.Tag(i)

		if content := g.processField(ctx, field, tag); content != nil {
			node.Content = append(node.Content, content...)
		}
	}

	return node
}

func (g *YAML) processField(
	ctx context.Context,
	field *types.Var,
	tag string,
) []*yaml.Node {
	g.debugf("processing field %s", field.Name())

	yamlName := gentype.ParseNameTag(tag, g.OutputOptions.Tag, field.Name())
	if yamlName == "" {
		g.debugf("skipping field %s", field.Name())

		return nil
	}

	value := gentype.ParseDefaultValue(tag, gentype.ValueTagsYAML()...)
	comment := g.Source.CommentsMap[field.Pos()]
	ft := field.Type()

	if field.Anonymous() {
		if stt, _, ok := gentype.GetUnderlyingStruct(ft); ok {
			embedded := g.structToYAMLNode(ctx, stt)
			if embedded == nil {
				return nil
			}

			return embedded.Content
		}
	}

	if !field.Exported() {
		return nil
	}

	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: yamlName}
	valNode := g.typeToYAMLNode(ctx, ft, value)

	if comment != "" {
		keyNode.HeadComment = comment
	}

	return []*yaml.Node{keyNode, valNode}
}

//nolint:cyclop,funlen
func (g *YAML) typeToYAMLNode(ctx context.Context, t types.Type, value string) *yaml.Node {
	if gentype.IsTextUnmarshaler(t) || gentype.IsStringer(t) {
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
		return g.typeToYAMLNode(ctx, tt.Elem(), value)
	case *types.Slice, *types.Array:
		elemType := tt.(interface{ Elem() types.Type }).Elem()

		seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

		if value != "" && !isStructLike(elemType) {
			for _, v := range strings.Split(value, ",") {
				elemNode := g.typeToYAMLNode(ctx, elemType, v)
				seq.Content = append(seq.Content, elemNode)
			}
		} else {
			seq.Content = append(seq.Content, g.typeToYAMLNode(ctx, elemType, ""))
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
			if _, ok := g.visited[key]; ok {
				return &yaml.Node{
					Kind:    yaml.MappingNode,
					Tag:     "!!map",
					Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "<recursive>"}},
				}
			}

			g.visited[key] = struct{}{}

			return g.structToYAMLNode(ctx, st)
		}

		return g.typeToYAMLNode(ctx, tt.Underlying(), value)
	case *types.Struct:
		key := tt.String()
		if _, ok := g.visited[key]; ok {
			return &yaml.Node{
				Kind:    yaml.MappingNode,
				Tag:     "!!map",
				Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "<recursive>"}},
			}
		}

		g.visited[key] = struct{}{}

		return g.structToYAMLNode(ctx, tt)
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: ""}
}
