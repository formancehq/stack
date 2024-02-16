package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/spf13/cobra"
)

func NewDatabaseCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "db",
		Short: "Handle databases operations",
	}
	ret.AddCommand(NewDatabaseCreateCommand(), NewDatabaseDropCommand())
	bunconnect.InitFlags(ret.PersistentFlags())

	return ret
}
