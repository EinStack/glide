package cmd

import (
	"github.com/spf13/cobra"
	"glide/pkg"
	"glide/pkg/telemetry"
	"go.uber.org/zap"
)

// NewCLI Create a Glide CLI
func NewCLI() *cobra.Command {
	// TODO: Chances are we could use the build in flags module in this is all we need from CLI
	cli := &cobra.Command{
		Use:     "",
		Version: pkg.GetVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: gonna be read from a config file
			logConfig := telemetry.NewLogConfig()
			logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			logConfig.Encoding = "console"
			logger, err := telemetry.NewLogger(logConfig)

			if err != nil {
				return err
			}

			logger.Debug("logger inited")

			gateway, err := pkg.NewGateway()

			if err != nil {
				return err
			}

			return gateway.Run(cmd.Context())
		},
		//SilenceUsage: true,
	}

	return cli
}
