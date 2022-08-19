package cmd

import (
	"github.com/numary/webhooks/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run webhooks worker",
	RunE:  worker.Run,
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
