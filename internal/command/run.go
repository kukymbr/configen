package command

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/kukymbr/configen/internal/formatter"
	"github.com/kukymbr/configen/internal/generator"
	"github.com/kukymbr/configen/internal/logger"
	"github.com/kukymbr/configen/internal/version"
)

func Run() {
	if err := run(); err != nil {
		logger.Errorf("%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	opt := generator.Options{}
	silent := false

	var cmd = &cobra.Command{
		Use:   "configen",
		Short: "Golang code generator",
		Long:  `Golang code generator`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			gen, err := generator.New(opt)
			if err != nil {
				return err
			}

			return gen.Generate(ctx)
		},
		Version: version.GetVersion(),
	}

	initFlags(cmd, &opt, &silent)

	cmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		logger.SetSilentMode(silent)
	}

	return cmd.Execute()
}

func initFlags(cmd *cobra.Command, opt *generator.Options, silent *bool) {
	cmd.PersistentFlags().BoolVarP(silent, "silent", "s", false, "Silent mode")

	cmd.Flags().StringVar(
		&opt.PackageName,
		"package",
		generator.DefaultPackageName,
		"Target package name of the generated code",
	)

	cmd.Flags().StringVar(
		&opt.TargetDir,
		"target",
		generator.DefaultTargetDir,
		"Directory for the generated Go files",
	)

	cmd.Flags().StringVar(
		&opt.Formatter,
		"fmt",
		generator.DefaultFormatter,
		"Formatter used to format generated go files ("+formatter.GoFmt+"|"+formatter.Noop+")",
	)
}
