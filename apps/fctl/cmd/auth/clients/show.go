package clients

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/auth/clients/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Client *shared.Client `json:"client,omitempty"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewDefaultShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <client-id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Show client"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	authClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.ReadClientRequest{
		ClientID: args[0],
	}
	response, err := authClient.Auth.ReadClient(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Client = response.ReadClientResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {

	views.PrintClient(cmd.OutOrStdout(), c.store.Client)

	if len(c.store.Client.RedirectUris) > 0 {
		fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Redirect URIs :")
		if err := pterm.DefaultBulletList.WithWriter(cmd.OutOrStdout()).WithItems(fctl.Map(c.store.Client.RedirectUris, func(redirectURI string) pterm.BulletListItem {
			return pterm.BulletListItem{
				Text:        redirectURI,
				TextStyle:   pterm.NewStyle(pterm.FgDefault),
				BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
			}
		})).Render(); err != nil {
			return err
		}
	}

	if len(c.store.Client.PostLogoutRedirectUris) > 0 {
		fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Post logout redirect URIs :")
		if err := pterm.DefaultBulletList.WithWriter(cmd.OutOrStdout()).WithItems(fctl.Map(c.store.Client.PostLogoutRedirectUris, func(redirectURI string) pterm.BulletListItem {
			return pterm.BulletListItem{
				Text:        redirectURI,
				TextStyle:   pterm.NewStyle(pterm.FgDefault),
				BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
			}
		})).Render(); err != nil {
			return err
		}
	}

	if len(c.store.Client.Secrets) > 0 {
		if err := views.PrintSecrets(cmd.OutOrStdout(), c.store.Client.Secrets); err != nil {
			return err
		}
	}
	return nil
}
