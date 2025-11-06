package env

import (
	"context"
	"fmt"
	"go/types"
	"reflect"

	"github.com/kukymbr/configen/internal/generator/gentype"
)

func (g *Env) collectEnvVars(ctx context.Context, st *types.Struct, prefix string) {
	ctx = gentype.ContextIncRecursionDepth(ctx)
	gentype.ContextMustValidateRecursionDepth(ctx, "Env generator (collectEnvVars)")

	if ctx.Err() != nil {
		return
	}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		tag := st.Tag(i)

		g.processField(ctx, field, tag, prefix)
	}
}

//nolint:cyclop
func (g *Env) processField(ctx context.Context, field *types.Var, tag string, prefix string) {
	envName := gentype.ParseNameTag(tag, g.OutputOptions.Tag, "")
	envPrefix := reflect.StructTag(tag).Get(gentype.TagEnvPrefix)
	example := gentype.ParseDefaultValue(tag, gentype.ValueTagsEnv(g.OutputOptions.DefaultValueTag)...)
	value := gentype.NewNullable[string]()

	comment := g.Source.CommentsMap[field.Pos()]
	ft := field.Type()

	if field.Anonymous() {
		g.processAnonymousField(ctx, ft, prefix)

		return
	}

	if !field.Exported() {
		return
	}

	if stt, named, ok := gentype.GetUnderlyingStruct(ft); ok {
		if len(g.envs) > 0 && g.envs[len(g.envs)-1] != "" {
			// Separate substructs with space.
			g.envs = append(g.envs, "")
		}

		if named == nil {
			g.collectEnvVars(ctx, stt, prefix+envPrefix)

			return
		}

		if !g.isTargetPackage(named) {
			return
		}

		if gentype.IsTextUnmarshaler(stt) {
			value.Set(gentype.DefaultValueForType(ft, example))
		} else {
			g.collectEnvVars(ctx, stt, prefix+envPrefix)
		}
	}

	if envName == "" {
		return
	}

	if !value.IsSet() {
		value.Set(gentype.DefaultValueForType(ft, example))
	}

	if comment != "" {
		g.envs = append(g.envs, fmt.Sprintf("# %s", comment))
	}

	g.envs = append(g.envs, fmt.Sprintf("%s%s=%s", prefix, envName, value.Value()))
}

// processAnonymousField expands anonymous embedded struct fields in env values.
func (g *Env) processAnonymousField(ctx context.Context, ft types.Type, prefix string) {
	stt, named, ok := gentype.GetUnderlyingStruct(ft)
	if !ok {
		return
	}

	if named == nil {
		return
	}

	g.collectEnvVars(ctx, stt, prefix)
}

func (g *Env) isTargetPackage(name any) bool {
	return gentype.ParsePackageName(name) == gentype.ParsePackageName(g.Source.Package)
}
