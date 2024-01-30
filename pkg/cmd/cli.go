package cmd

import (
	"glide/pkg"
	"glide/pkg/config"

	"github.com/spf13/cobra"
)

var cfgFile string

const Description = `
 ██████╗ ██╗     ██╗██████╗ ███████╗
██╔════╝ ██║     ██║██╔══██╗██╔════╝
██║  ███╗██║     ██║██║  ██║█████╗  
██║   ██║██║     ██║██║  ██║██╔══╝  
╚██████╔╝███████╗██║██████╔╝███████╗
 ╚═════╝ ╚══════╝╚═╝╚═════╝ ╚══════╝
🐦An open-source, lightweight, high-performance model gateway 
to make your LLM applications production ready 🎉

📚Documentation: https://glide.einstack.ai
🛠️Source: https://github.com/EinStack/glide
💬Discord: https://discord.gg/pt53Ej7rrc
🐛Bug Tracker: https://github.com/EinStack/glide/issues

🏗️EinStack Community (mailto:contact@einstack.ai), 2024-Present (c)
`

// NewCLI Create a Glide CLI
func NewCLI() *cobra.Command {
	// TODO: Chances are we could use the build in flags module in this is all we need from CLI
	cli := &cobra.Command{
		Use:     "glide",
		Short:   "🐦Glide is an open-source, lightweight, high-performance model gateway",
		Long:    Description,
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
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cli.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	_ = cli.MarkPersistentFlagRequired("config")

	return cli
}
