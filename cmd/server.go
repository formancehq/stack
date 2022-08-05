package cmd

import (
	"github.com/numary/webhooks-cloud/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch server",
	RunE:  server.Start,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
