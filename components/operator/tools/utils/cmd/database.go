package cmd

import (
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/spf13/cobra"
)

func NewDatabaseCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "db",
		Short: "Handle databases operations",
	}
	ret.AddCommand(NewDatabaseCreateCommand(), NewDatabaseDropCommand())
	bunconnect.AddFlags(ret.PersistentFlags())

	return ret
}
