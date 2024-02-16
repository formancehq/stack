package cmd

import (
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/service"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "utils",
	Short: "A cli for operator operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := service.GetDefaultLogger(cmd.OutOrStdout())
		logger.Infof("Starting application")
		logger.Debugf("Environment variables:")
		for _, v := range os.Environ() {
			logger.Debugf(v)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewDatabaseCommand())
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	service.BindFlags(rootCmd)

	viper.SetEnvKeyReplacer(EnvVarReplacer)
	viper.AutomaticEnv()
}

var EnvVarReplacer = strings.NewReplacer(".", "_", "-", "_")
