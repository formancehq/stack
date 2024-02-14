package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDatabaseDropCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "drop",
		Short: "Handle database dropping",
		RunE: func(cmd *cobra.Command, args []string) error {
			connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.OutOrStdout(), viper.GetBool(service.DebugFlag))
			if err != nil {
				return err
			}

			return bunmigrate.EnsureDatabaseExists(cmd.Context(), *connectionOptions)
		},
	}
	return ret
}
