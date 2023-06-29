package clients

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <client-id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Show client"),
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

			authClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			request := operations.ReadClientRequest{
				ClientID: args[0],
			}
			response, err := authClient.Auth.ReadClient(cmd.Context(), request)
			if err != nil {
				return err
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), response.ReadClientResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Name"), response.ReadClientResponse.Data.Name})
			tableData = append(tableData, []string{pterm.LightCyan("Description"), fctl.StringPointerToString(response.ReadClientResponse.Data.Description)})
			tableData = append(tableData, []string{pterm.LightCyan("Public"), fctl.BoolPointerToString(response.ReadClientResponse.Data.Public)})

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information :")
			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "")

			if len(response.ReadClientResponse.Data.RedirectUris) > 0 {
				fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Redirect URIs :")
				if err := pterm.DefaultBulletList.WithWriter(cmd.OutOrStdout()).WithItems(fctl.Map(response.ReadClientResponse.Data.RedirectUris, func(redirectURI string) pterm.BulletListItem {
					return pterm.BulletListItem{
						Text:        redirectURI,
						TextStyle:   pterm.NewStyle(pterm.FgDefault),
						BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
					}
				})).Render(); err != nil {
					return err
				}
			}

			if len(response.ReadClientResponse.Data.PostLogoutRedirectUris) > 0 {
				fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Post logout redirect URIs :")
				if err := pterm.DefaultBulletList.WithWriter(cmd.OutOrStdout()).WithItems(fctl.Map(response.ReadClientResponse.Data.PostLogoutRedirectUris, func(redirectURI string) pterm.BulletListItem {
					return pterm.BulletListItem{
						Text:        redirectURI,
						TextStyle:   pterm.NewStyle(pterm.FgDefault),
						BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
					}
				})).Render(); err != nil {
					return err
				}
			}

			if len(response.ReadClientResponse.Data.Secrets) > 0 {
				fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Secrets :")

				if err := pterm.DefaultTable.
					WithWriter(cmd.OutOrStdout()).
					WithHasHeader(true).
					WithData(fctl.Prepend(
						fctl.Map(response.ReadClientResponse.Data.Secrets, func(secret shared.ClientSecret) []string {
							return []string{
								secret.ID, secret.Name, secret.LastDigits,
							}
						}),
						[]string{"ID", "Name", "Last digits"},
					)).
					Render(); err != nil {
					return err
				}
			}

			return nil
		}),
	)
}
