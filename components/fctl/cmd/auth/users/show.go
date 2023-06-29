package users

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <user-id>",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Show user"),
		fctl.WithArgs(cobra.ExactArgs(1)),
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

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			request := operations.ReadUserRequest{
				UserID: args[0],
			}
			readUserResponse, err := client.Auth.ReadUser(cmd.Context(), request)
			if err != nil {
				return err
			}

			if readUserResponse.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", readUserResponse.StatusCode)
			}

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), *readUserResponse.ReadUserResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Membership ID"), *readUserResponse.ReadUserResponse.Data.Subject})
			tableData = append(tableData, []string{pterm.LightCyan("Email"), *readUserResponse.ReadUserResponse.Data.Email})

			return pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
