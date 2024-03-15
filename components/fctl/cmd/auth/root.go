package auth

import (
	"github.com/formancehq/fctl/cmd/auth/clients"
	"github.com/formancehq/fctl/cmd/auth/users"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("auth",
		fctl.WithShortDescription("Auth server management"),
		fctl.WithChildCommands(
			clients.NewCommand(),
			users.NewCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
