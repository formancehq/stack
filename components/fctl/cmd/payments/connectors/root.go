package connectors

import (
	"github.com/formancehq/fctl/cmd/payments/connectors/configs"
	"github.com/formancehq/fctl/cmd/payments/connectors/install"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewConnectorsCommand() *cobra.Command {
	return fctl.NewCommand("connectors",
		fctl.WithAliases("c", "co", "con"),
		fctl.WithShortDescription("Connectors management"),
		fctl.WithChildCommands(
			NewUninstallCommand(),
			NewListCommand(),
			install.NewInstallCommand(),
			configs.NewUpdateConfigCommands(),
			configs.NewGetConfigCommand(),
		),
	)
}
