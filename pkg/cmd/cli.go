package cmd

import (
	"github.com/EinStack/glide/pkg/telemetry"
	"github.com/EinStack/glide/pkg/version"
	"go.uber.org/zap"

	"github.com/EinStack/glide/pkg/config"

	"github.com/EinStack/glide/pkg"

	"github.com/spf13/cobra"
)

var (
	dotEnvFile string
	cfgFile    string
	logger     = telemetry.GetLogger()
)

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
		Use:     "github.com/EinStack/glide",
		Short:   "🐦Glide is an open-source, lightweight, high-performance model gateway",
		Long:    Description,
		Version: version.FullVersion,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configProvider := config.NewProvider()

			err := configProvider.LoadDotEnv(dotEnvFile)

			if err != nil {
				logger.Warn("⚠️failed to load dotenv file: ", zap.Error(err)) // don't have an inited logger at this moment
			} else {
				logger.Info("🔧dot env file is loaded", zap.Any("dotenv", dotEnvFile))
			}

			_, err = configProvider.Load(cfgFile)
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

	cli.PersistentFlags().StringVarP(&dotEnvFile, "env", "e", ".env", "dotenv file")
	cli.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	_ = cli.MarkPersistentFlagRequired("config")

	return cli
}
