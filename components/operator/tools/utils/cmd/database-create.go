package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewDatabaseCreateCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "create",
		Short: "Handle database creation",
		RunE: func(cmd *cobra.Command, args []string) error {
			connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
			if err != nil {
				return errors.Wrap(err, "resolving connection params")
			}

			return errors.Wrap(
				bunmigrate.EnsureDatabaseExists(cmd.Context(), *connectionOptions),
				"ensuring database exists",
			)
		},
	}
	return ret
}
