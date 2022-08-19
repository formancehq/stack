package cmd

import (
	"github.com/numary/webhooks/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run webhooks server",
	RunE:  server.Run,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
