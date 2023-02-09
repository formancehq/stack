package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/formancehq/stack/components/gateway/pkg/builder"
)

const (
	caddyBuilderConfigPath = "caddy-builder-config-path"
	caddyBinaryOutputPath  = "caddy-binary-output-path"
)

var buildCmd = &cobra.Command{
	Use: "build",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := builder.NewBuilder(
			viper.GetString(caddyBuilderConfigPath),
			viper.GetString(caddyBinaryOutputPath),
		)

		return builder.Build(cmd.Context())
	},
}

func init() {
	buildCmd.Flags().String(caddyBuilderConfigPath, "", "Path to the caddy builder's config file")
	buildCmd.MarkFlagRequired(caddyBuilderConfigPath)
	buildCmd.Flags().String(caddyBinaryOutputPath, "", "Path to the caddy binary output file")
	buildCmd.MarkFlagRequired(caddyBinaryOutputPath)
	rootCmd.AddCommand(buildCmd)
}
