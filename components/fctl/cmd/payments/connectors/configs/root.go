package configs

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewUpdateConfigCommands() *cobra.Command {
	return fctl.NewCommand("update-config",
		fctl.WithAliases("uc"),
		fctl.WithShortDescription("Update the config of a connector"),
		fctl.WithChildCommands(
			newUpdateAdyenCommand(),
			newUpdateAtlarCommand(),
			newUpdateBankingCircleCommand(),
			newUpdateCurrencyCloudCommand(),
			newUpdateMangopayCommand(),
			newUpdateModulrCommand(),
			newUpdateMoneycorpCommand(),
			newUpdateStripeCommand(),
			newUpdateWiseCommand(),
		),
	)
}
