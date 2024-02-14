package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "utils",
	Short: "A cli for operator operations",
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

	viper.SetEnvKeyReplacer(EnvVarReplacer)
	viper.AutomaticEnv()
}

var EnvVarReplacer = strings.NewReplacer(".", "_", "-", "_")
