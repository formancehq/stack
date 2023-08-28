package clients

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/auth/clients/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show <client-id>"
	shortShow = "Show client details"
)

type ShowStore struct {
	Client *shared.Client `json:"client,omitempty"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.ReadClientRequest{
		ClientID: args[0],
	}
	response, err := authClient.Auth.ReadClient(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Client = response.ReadClientResponse.Data

	return c, nil
}

func (c *ShowController) Render() error {
	out := c.config.GetOut()
	err := views.PrintClient(out, c.store.Client)
	if err != nil {
		return errors.Wrap(err, "failed to print client")
	}

	if len(c.store.Client.RedirectUris) > 0 {
		fctl.BasicTextCyan.WithWriter(out).Printfln("Redirect URIs :")
		if err := pterm.DefaultBulletList.WithWriter(out).WithItems(fctl.Map(c.store.Client.RedirectUris, func(redirectURI string) pterm.BulletListItem {
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
		fctl.BasicTextCyan.WithWriter(out).Printfln("Post logout redirect URIs :")
		if err := pterm.DefaultBulletList.WithWriter(out).WithItems(fctl.Map(c.store.Client.PostLogoutRedirectUris, func(redirectURI string) pterm.BulletListItem {
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
		if err := views.PrintSecrets(out, c.store.Client.Secrets); err != nil {
			return err
		}
	}
	return nil
}
func NewShowCommand() *cobra.Command {
	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
