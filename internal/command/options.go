package command

import (
	"strings"

	"github.com/kukymbr/configen/internal/generator"
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

	// SourceDir is a directory of the source go files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string
}

func (opt options) ToGeneratorOptions() generator.Options {
	gen := generator.Options{
		StructName: opt.StructName,
		SourceDir:  opt.SourceDir,
	}

	outOpts := []struct {
		Input  string
		Target *generator.OutputOptions
	}{
		{Input: opt.YAMLPath, Target: &gen.YAML},
		{Input: opt.EnvPath, Target: &gen.Env},
	}

	for _, out := range outOpts {
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

	return gen
}
