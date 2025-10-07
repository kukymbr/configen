package command

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kukymbr/configen/internal/generator"
	"github.com/kukymbr/configen/internal/logger"
	"github.com/kukymbr/configen/internal/version"
	"github.com/spf13/cobra"
)

func Run() error {
	opt := options{}
	silent := false

	var cmd = &cobra.Command{
		Use:   "configen",
		Short: "Configs generator",
		Long:  `The go:generate tool to generate YAML and dotenv configuration files from the Golang struct.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			gen, err := generator.New(opt.ToGeneratorOptions())
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

func initFlags(cmd *cobra.Command, opt *options, silent *bool) {
	cmd.PersistentFlags().BoolVarP(silent, "silent", "s", false, "Silent mode")

	cmd.Flags().StringVar(
		&opt.StructName,
		"struct",
		"",
		"Name of the struct to generate config from",
	)

	cmd.Flags().StringVar(
		&opt.YAMLPath,
		"yaml",
		"",
		"Path to YAML config file, set 'true' to enable with default path",
	)

	cmd.Flags().StringVar(
		&opt.EnvPath,
		"env",
		"",
		"Path to dotenv config file, set 'true' to enable with default path",
	)

	cmd.Flags().StringVar(
		&opt.SourceDir,
		"source",
		generator.DefaultSourceDir,
		"Directory of the source go files",
	)

	_ = cmd.MarkFlagRequired("struct")
	cmd.MarkFlagsOneRequired("yaml", "env")
	_ = cmd.MarkFlagFilename("yaml")
	_ = cmd.MarkFlagFilename("env")
	_ = cmd.MarkFlagDirname("source")
}
