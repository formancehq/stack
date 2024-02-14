package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDatabaseCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "db",
		Short: "Handle databases operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.PersistentFlags())
		},
	}
	ret.AddCommand(NewDatabaseCreateCommand(), NewDatabaseDropCommand())
	bunconnect.InitFlags(ret.PersistentFlags())
	return ret
}
