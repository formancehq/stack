package clients

import (
	"flag"
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/auth/clients/internal"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// TODO: This command is a copy/paste of the create command
// We should handle membership side the patch of the client OR
// We should get the client before updating it to get replace informations

const (
	useUpdate   = "update <client-id>"
	shortUpdate = "Update a client"
)

type UpdateClient struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	IsPublic              string `json:"isPublic"`
	RedirectUri           string `json:"redirectUri"`
	PostLogoutRedirectUri string `json:"postLogoutRedirectUri"`
}

type UpdateStore struct {
	Client *UpdateClient `json:"client"`
}

func NewUpdateStore() *UpdateStore {
	return &UpdateStore{
		Client: &UpdateClient{},
	}
}

func NewUpdateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useUpdate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	flags.Bool(internal.PublicFlag, false, "Is client public")
	flags.Bool(internal.TrustedFlag, false, "Is the client trusted")
	flags.String(internal.DescriptionFlag, "", "Client description")
	flags.String(internal.RedirectUriFlag, "", "Redirect URIS")
	flags.String(internal.PostLogoutRedirectUriFlag, "", "Post logout redirect uris")
	return fctl.NewControllerConfig(
		useUpdate,
		shortUpdate,
		shortUpdate,
		[]string{
			"u", "upd",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

type UpdateController struct {
	store  *UpdateStore
	config *fctl.ControllerConfig
}

func NewUpdateController(config *fctl.ControllerConfig) *UpdateController {
	return &UpdateController{
		store:  NewUpdateStore(),
		config: config,
	}
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *UpdateController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete an OAuth2 client") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	public := fctl.GetBool(flags, internal.PublicFlag)
	trusted := fctl.GetBool(flags, internal.TrustedFlag)
	description := fctl.GetString(flags, internal.DescriptionFlag)

	request := operations.UpdateClientRequest{
		ClientID: args[0],
		UpdateClientRequest: &shared.UpdateClientRequest{
			Public:                 &public,
			RedirectUris:           fctl.GetStringSlice(flags, internal.RedirectUriFlag),
			Description:            &description,
			Name:                   args[0],
			Trusted:                &trusted,
			PostLogoutRedirectUris: fctl.GetStringSlice(flags, internal.PostLogoutRedirectUriFlag),
		},
	}
	response, err := authClient.Auth.UpdateClient(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Client.ID = response.UpdateClientResponse.Data.ID
	c.store.Client.Name = response.UpdateClientResponse.Data.Name
	c.store.Client.Description = fctl.StringPointerToString(response.UpdateClientResponse.Data.Description)
	c.store.Client.IsPublic = fctl.BoolPointerToString(response.UpdateClientResponse.Data.Public)
	c.store.Client.RedirectUri = strings.Join(response.UpdateClientResponse.Data.RedirectUris, ",")
	c.store.Client.PostLogoutRedirectUri = strings.Join(response.UpdateClientResponse.Data.PostLogoutRedirectUris, ",")

	return c, nil
}

func (c *UpdateController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Client.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Client.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), c.store.Client.Description})
	tableData = append(tableData, []string{pterm.LightCyan("Public"), c.store.Client.IsPublic})
	tableData = append(tableData, []string{pterm.LightCyan("Redirect URIs"), c.store.Client.RedirectUri})
	tableData = append(tableData, []string{pterm.LightCyan("Post logout redirect URIs"), c.store.Client.PostLogoutRedirectUri})
	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewUpdateCommand() *cobra.Command {
	config := NewUpdateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*UpdateStore](NewUpdateController(config)),
	)
}
