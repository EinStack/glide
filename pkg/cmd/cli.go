package cmd

import (
	"github.com/spf13/cobra"
	"glide/pkg"
)

// NewCLI Create a Glide CLI
func NewCLI() *cobra.Command {
	// TODO: Chances are we could use the build in flags module in this is all we need from CLI
	cli := &cobra.Command{
		Use:     "",
		Version: pkg.GetVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			gateway, err := pkg.NewGateway()
			if err != nil {
				return err
			}

			return gateway.Run(cmd.Context())
		},
		// SilenceUsage: true,
	}

	return cli
}
