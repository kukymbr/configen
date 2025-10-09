package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kukymbr/configen/internal/generator/gentype"
)

const (
	DefaultSourceDir          = "."
	DefaultEnvTag             = "env"
	DefaultYAMLTag            = "yaml"
	DefaultGoTargetStructName = "ConfigProvider"
)

type Options struct {
	// StructName is a struct name to generate config from.
	StructName string

	// YAML target config file options
	YAML gentype.OutputOptions

	// Env target config file options
	Env gentype.OutputOptions

	GoGetter gentype.OutputOptions

	// SourceDir is a directory of the SQL files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string
}

func (opt Options) Debug() string {
	return fmt.Sprintf("%#v", opt)
}

//nolint:cyclop
func prepareOptions(opt *Options) error {
	if opt.StructName == "" {
		return fmt.Errorf("struct name is required")
	}

	if err := validateIdentifier(opt.StructName); err != nil {
		return err
	}

	if opt.SourceDir == "" {
		opt.SourceDir = DefaultSourceDir
	}

	structSlug := strings.ToLower(opt.StructName)

	if opt.YAML.Path == "" {
		opt.YAML.Path = structSlug + ".yaml"
	}

	if opt.Env.Path == "" {
		opt.Env.Path = structSlug + ".env"
	}

	if opt.GoGetter.Path == "" {
		opt.GoGetter.Path = structSlug + ".go"
	}

	if opt.YAML.Tag == "" {
		opt.YAML.Tag = DefaultYAMLTag
	}

	if opt.Env.Tag == "" {
		opt.Env.Tag = DefaultEnvTag
	}

	if opt.GoGetter.TargetStructName == "" {
		opt.GoGetter.TargetStructName = DefaultGoTargetStructName
	}

	if err := ensureDirs(opt.YAML, opt.Env, opt.GoGetter); err != nil {
		return err
	}

	if err := validateIsDir(opt.SourceDir); err != nil {
		return err
	}

	return nil
}

func ensureDirs(opts ...gentype.OutputOptions) error {
	for _, opt := range opts {
		if dir := filepath.Dir(opt.Path); dir != "" {
			if err := EnsureDir(dir); err != nil {
				return err
			}
		}
	}

	return nil
}
