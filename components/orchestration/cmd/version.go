package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s \n", Version)
			fmt.Printf("Date: %s \n", BuildDate)
			fmt.Printf("Commit: %s \n", Commit)
		},
	}
}
