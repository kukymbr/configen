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
	// Set "true" for enable the generator with a default file path.
	EnvPath string

	GoPath string

	// YAMLTag is a tag name for a YAML field names, `yaml` by default.
	YAMLTag string

	// EnvTag is a tag name for a dotenv field names, `env` by default.
	EnvTag string

	// SourceDir is a directory of the source go files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string

	GoTargetStructName  string
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
	}

	gen.GoGetter.TargetStructName = opt.GoTargetStructName
	gen.GoGetter.TargetPackageName = opt.GoTargetPackageName

	return gen
}
