package secrets

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <client-id> <secret-name>",
		fctl.WithAliases("c"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithShortDescription("Create secret"),
		fctl.WithConfirmFlag(),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to create a new client secret") {
				return fctl.ErrMissingApproval
			}

			authClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			request := operations.CreateSecretRequest{
				ClientID: args[0],
				CreateSecretRequest: &shared.CreateSecretRequest{
					Name:     args[1],
					Metadata: nil,
				},
			}
			response, err := authClient.Auth.CreateSecret(cmd.Context(), request)
			if err != nil {
				return err
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), response.CreateSecretResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Name"), response.CreateSecretResponse.Data.Name})
			tableData = append(tableData, []string{pterm.LightCyan("Clear"), response.CreateSecretResponse.Data.Clear})
			return pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
