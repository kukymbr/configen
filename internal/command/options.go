package command

import (
	"strings"

	"github.com/kukymbr/configen/internal/generator"
	"github.com/kukymbr/configen/internal/generator/gentype"
)

const (
	keywordTrue  = "true"
	keywordFalse = "false"
)

type options struct {
	// StructName is a struct name to generate config from.
	StructName string

	// YAMLPath is a path to a target YAML config file.
	// Define to enable YAML generator.
	// Set "true" for enable the generator with a default file path.
	YAMLPath string

	// EnvPath is a path to a target .env config file.
	// Define to enable Env generator.
	// Set "true" to enable the generator with a default file path.
	EnvPath string

	// GoPath is a path to a target Go config getter file.
	GoPath string

	// YAMLTag is a tag name for YAML field name, `yaml` by default.
	YAMLTag string

	// EnvTag is a tag name for a dotenv field name, `env` by default.
	EnvTag string

	// DefaultValueTag is an explicit tag name for a default value.
	// Overrides the default lookup if given.
	DefaultValueTag string

	// SourceDir is a directory of the source go files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string

	// GoTargetStructName is the name of the target struct.
	GoTargetStructName string

	// GoTargetPackageName is the name of the target golang package.
	GoTargetPackageName string
}

func (opt options) ToGeneratorOptions() generator.Options {
	gen := generator.Options{
		StructName: opt.StructName,
		SourceDir:  opt.SourceDir,
	}

	outOpts := []struct {
		Input  string
		Tag    string
		Target *gentype.OutputOptions
	}{
		{Input: opt.YAMLPath, Tag: opt.YAMLTag, Target: &gen.YAML},
		{Input: opt.EnvPath, Tag: opt.EnvTag, Target: &gen.Env},
		{Input: opt.GoPath, Target: &gen.GoGetter},
	}

	for _, out := range outOpts {
		out.Target.Tag = out.Tag

		keyword := strings.ToLower(out.Input)
		switch keyword {
		case keywordTrue:
			out.Target.Enable = true

			continue
		case keywordFalse, "":
			continue
		}

		out.Target.Enable = true
		out.Target.Path = out.Input
		out.Target.DefaultValueTag = opt.DefaultValueTag
	}

	gen.GoGetter.TargetStructName = opt.GoTargetStructName
	gen.GoGetter.TargetPackageName = opt.GoTargetPackageName

	return gen
}
