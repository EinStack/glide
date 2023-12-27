package cmd

import (
	"glide/pkg"
	"glide/pkg/config"

	"github.com/spf13/cobra"
)

var cfgFile string

// NewCLI Create a Glide CLI
func NewCLI() *cobra.Command {
	// TODO: Chances are we could use the build in flags module in this is all we need from CLI
	cli := &cobra.Command{
		Use:     "glide",
		Short:   "🐦Glide is an open-source, lightweight, high-performance model gateway",
		Long:    "TODO",
		Version: pkg.FullVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			configProvider, err := config.NewProvider().Load(cfgFile)
			if err != nil {
				return err
			}

			gateway, err := pkg.NewGateway(configProvider)
			if err != nil {
				return err
			}

			return gateway.Run(cmd.Context())
		},
		// SilenceUsage: true,
	}

	cli.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	return cli
}
