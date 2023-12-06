package cmd

import "github.com/spf13/cobra"

func newMigrate() *cobra.Command {
	migrate := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE:  runMigrate,
	}

	return migrate
}

func runMigrate(cmd *cobra.Command, args []string) error {
	return nil
}
