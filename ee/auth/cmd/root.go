package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/uptrace/bun"

	"github.com/spf13/cobra"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	cobra.EnableTraverseRunHooks = true

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.AddCommand(
		newServeCommand(),
		newVersionCommand(),
		bunmigrate.NewDefaultCommand(func(cmd *cobra.Command, args []string, db *bun.DB) error {
			return sqlstorage.Migrate(cmd.Context(), db)
		}))

	return cmd
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		exitWithCode(1, err)
	}
}
