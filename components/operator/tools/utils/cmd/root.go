package cmd

import (
	"fmt"
	"os"
	"strings"

	logging "github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/formancehq/stack/libs/go-libs/service"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "utils",
	Short: "A cli for operator operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := service.GetDefaultLogger(cmd.OutOrStdout(), viper.GetBool(service.DebugFlag))
		logger.Infof("Starting application")
		logger.Debugf("Environment variables:")
		for _, v := range os.Environ() {
			logger.Debugf(v)
		}
		cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
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
	if err := bindFlagsToViper(rootCmd); err != nil {
		panic(err)
	}

	viper.SetEnvKeyReplacer(EnvVarReplacer)
	viper.AutomaticEnv()
}

func bindFlagsToViper(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	return viper.BindPFlags(cmd.PersistentFlags())
}

var EnvVarReplacer = strings.NewReplacer(".", "_", "-", "_")
