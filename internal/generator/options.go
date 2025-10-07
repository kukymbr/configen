package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kukymbr/configen/internal/utils"
)

const (
	DefaultSourceDir = "."
	DefaultEnvTag    = "env"
	DefaultYAMLTag   = "yaml"
)

type Options struct {
	// StructName is a struct name to generate config from.
	StructName string

	// YAML target config file options
	YAML OutputOptions

	// Env target config file options
	Env OutputOptions

	// SourceDir is a directory of the SQL files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string
}

type OutputOptions struct {
	// Enable is a flag to enable an output.
	Enable bool

	// Path is a target file path.
	Path string

	// Tag is a field names tag.
	Tag string
}

func (opt Options) Debug() string {
	return fmt.Sprintf("%#v", opt)
}

func prepareOptions(opt *Options) error {
	if opt.StructName == "" {
		return fmt.Errorf("struct name is required")
	}

	if err := utils.ValidateIdentifier(opt.StructName); err != nil {
		return err
	}

	if opt.SourceDir == "" {
		opt.SourceDir = DefaultSourceDir
	}

	if opt.YAML.Path == "" {
		opt.YAML.Path = strings.ToLower(opt.StructName) + ".yaml"
	}

	if opt.Env.Path == "" {
		opt.Env.Path = strings.ToLower(opt.StructName) + ".env"
	}

	if opt.YAML.Tag == "" {
		opt.YAML.Tag = DefaultYAMLTag
	}

	if opt.Env.Tag == "" {
		opt.Env.Tag = DefaultEnvTag
	}

	if err := ensureDirs(opt.YAML, opt.Env); err != nil {
		return err
	}

	if err := utils.ValidateIsDir(opt.SourceDir); err != nil {
		return err
	}

	return nil
}

func ensureDirs(opts ...OutputOptions) error {
	for _, opt := range opts {
		if dir := filepath.Dir(opt.Path); dir != "" {
			if err := utils.EnsureDir(dir); err != nil {
				return err
			}
		}
	}

	return nil
}
