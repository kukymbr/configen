package generator

import (
	"fmt"
	"strings"

	"github.com/kukymbr/configen/internal/formatter"
	"github.com/kukymbr/configen/internal/utils"
)

const (
	DefaultPackageName = "configen"
	DefaultSourceDir   = "."
	DefaultTargetDir   = "."
	DefaultFormatter   = formatter.GoFmt
)

type Options struct {
	// PackageName is a target package name of the generated code.
	// Default is "configen".
	PackageName string

	// SourceDir is a directory of the SQL files.
	// Default is the current directory (most applicable for go:generate).
	SourceDir string

	// TargetDir is a target go code directory.
	// Default is the current directory.
	TargetDir string

	// Formatter is a name of the formatter for the generated code files.
	// Available options: gofmt (default), none.
	Formatter string
}

func (opt Options) Debug() string {
	return fmt.Sprintf("%#v", opt)
}

func prepareOptions(opt *Options) error {
	opt.PackageName = strings.TrimSpace(opt.PackageName)

	if opt.PackageName == "" {
		opt.PackageName = DefaultPackageName
	}

	if opt.SourceDir == "" {
		opt.SourceDir = DefaultSourceDir
	}

	if opt.TargetDir == "" {
		opt.TargetDir = DefaultTargetDir
	}

	if opt.Formatter == "" {
		opt.Formatter = DefaultFormatter
	}

	if err := utils.ValidateIsDir(opt.SourceDir); err != nil {
		return err
	}

	if err := utils.ValidatePackageName(opt.PackageName); err != nil {
		return err
	}

	if err := utils.EnsureDir(opt.TargetDir); err != nil {
		return err
	}

	return nil
}
