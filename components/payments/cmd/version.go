package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func newVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get version",
		Run:   printVersion,
	}
}

func printVersion(*cobra.Command, []string) {
	log.Printf("Version: %s \n", Version)
	log.Printf("Date: %s \n", BuildDate)
	log.Printf("Commit: %s \n", Commit)
}
