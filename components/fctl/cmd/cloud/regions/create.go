package regions

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create [name]",
		fctl.WithAliases("sh", "s"),
		fctl.WithShortDescription("Show region details"),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			name := ""
			if len(args) > 0 {
				name = args[0]
			} else {
				name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
				if err != nil {
					return err
				}
			}

			regionResponse, _, err := apiClient.DefaultApi.CreatePrivateRegion(cmd.Context(), organizationID).
				CreatePrivateRegionRequest(membershipclient.CreatePrivateRegionRequest{
					Name: name,
				}).
				Execute()
			if err != nil {
				return err
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
				"Region created successfully with ID: %s", regionResponse.Data.Id)
			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
				"Your secret is (keep it safe, we will not be able to give it to you again): %s", *regionResponse.Data.Secret.Clear)

			return nil
		}),
	)
}
