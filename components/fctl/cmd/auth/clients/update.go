package clients

import (
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// TODO: This command is a copy/paste of the create command
// We should handle membership side the patch of the client OR
// We should get the client before updating it to get replace informations
func NewUpdateCommand() *cobra.Command {
	const (
		publicFlag                = "public"
		trustedFlag               = "trusted"
		descriptionFlag           = "description"
		redirectUriFlag           = "redirect-uri"
		postLogoutRedirectUriFlag = "post-logout-redirect-uri"
	)
	return fctl.NewCommand("update <client-id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Update client"),
		fctl.WithAliases("u", "upd"),
		fctl.WithConfirmFlag(),
		fctl.WithBoolFlag(publicFlag, false, "Is client public"),
		fctl.WithBoolFlag(trustedFlag, false, "Is the client trusted"),
		fctl.WithStringFlag(descriptionFlag, "", "Client description"),
		fctl.WithStringSliceFlag(redirectUriFlag, []string{}, "Redirect URIS"),
		fctl.WithStringSliceFlag(postLogoutRedirectUriFlag, []string{}, "Post logout redirect uris"),
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

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete an OAuth2 client") {
				return fctl.ErrMissingApproval
			}

			authClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			public := fctl.GetBool(cmd, publicFlag)
			trusted := fctl.GetBool(cmd, trustedFlag)
			description := fctl.GetString(cmd, descriptionFlag)

			request := operations.UpdateClientRequest{
				ClientID: args[0],
				UpdateClientRequest: &shared.UpdateClientRequest{
					Public:                 &public,
					RedirectUris:           fctl.GetStringSlice(cmd, redirectUriFlag),
					Description:            &description,
					Name:                   args[0],
					Trusted:                &trusted,
					PostLogoutRedirectUris: fctl.GetStringSlice(cmd, postLogoutRedirectUriFlag),
				},
			}
			response, err := authClient.Auth.UpdateClient(cmd.Context(), request)
			if err != nil {
				return err
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), response.UpdateClientResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Name"), response.UpdateClientResponse.Data.Name})
			tableData = append(tableData, []string{pterm.LightCyan("Description"), fctl.StringPointerToString(response.UpdateClientResponse.Data.Description)})
			tableData = append(tableData, []string{pterm.LightCyan("Public"), fctl.BoolPointerToString(response.UpdateClientResponse.Data.Public)})
			tableData = append(tableData, []string{pterm.LightCyan("Redirect URIs"), strings.Join(response.UpdateClientResponse.Data.RedirectUris, ",")})
			tableData = append(tableData, []string{pterm.LightCyan("Post logout redirect URIs"), strings.Join(response.UpdateClientResponse.Data.PostLogoutRedirectUris, ",")})
			return pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
