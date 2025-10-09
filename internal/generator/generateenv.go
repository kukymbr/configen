package generator

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/kukymbr/configen/internal/generator/utils"
	"github.com/kukymbr/configen/internal/logger"
)

//nolint:unused
func generateEnv(src *gentype.Source, out gentype.OutputOptions) error {
	var envLines []string

	collectEnvVars(src.Struct, src.CommentsMap, "", &envLines, out.Tag)

	doc := gentype.GetDocComment("#", src.RootStructName, src.RootStructDoc)

	envContent := doc + strings.Join(envLines, "\n") + "\n"

	if err := utils.WriteFile([]byte(envContent), out.Path); err != nil {
		return err
	}

	logger.Successf("Generated %s file", out.Path)

	return nil
}

//nolint:unused
func collectEnvVars(
	st *types.Struct,
	comments map[token.Pos]string,
	prefix string,
	envs *[]string,
	tagName string,
) {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Exported() {
			continue
		}

		tag := st.Tag(i)
		envPrefix := reflect.StructTag(tag).Get(gentype.TagEnvPrefix)

		envName := gentype.ParseNameTag(tag, tagName, "")
		example := gentype.ParseDefaultValue(tag, gentype.ValueTagsEnv()...)
		comment := comments[field.Pos()]
		ft := field.Type()

		if stt, ok := gentype.GetUnderlyingStruct(ft); ok {
			if len(*envs) > 0 && (*envs)[len(*envs)-1] != "" {
				*envs = append(*envs, "")
			}

			collectEnvVars(stt, comments, prefix+envPrefix, envs, tagName)

			continue
		}

		if envName == "" {
			continue
		}

		val := gentype.DefaultValueForType(ft, example)

		if comment != "" {
			*envs = append(*envs, fmt.Sprintf("# %s", comment))
		}

		*envs = append(*envs, fmt.Sprintf("%s%s=%s", prefix, envName, val))
	}
}
