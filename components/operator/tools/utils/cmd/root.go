package cmd

import (
	"fmt"
	logging "github.com/formancehq/go-libs/logging"
	"os"

	"github.com/formancehq/go-libs/service"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "utils",
	Short: "A cli for operator operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.NewDefaultLogger(cmd.OutOrStdout(), service.IsDebug(cmd), false)
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

	service.AddFlags(rootCmd.PersistentFlags())
}
