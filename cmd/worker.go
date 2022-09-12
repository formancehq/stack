package cmd

import (
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use: "worker",
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
