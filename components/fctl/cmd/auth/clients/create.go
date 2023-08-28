package clients

import (
	"flag"
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/auth/clients/internal"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreate         = "create <name>"
	descriptionCreate = "Create a new client with the given name, and configure it with flags"
	shortCreate       = "Create a new client"
)

type CreateClient struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	IsPublic              string `json:"isPublic"`
	RedirectUri           string `json:"redirectUri"`
	PostLogoutRedirectUri string `json:"postLogoutRedirectUri"`
}

type CreateStore struct {
	Client *CreateClient `json:"client"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	flags.Bool(internal.PublicFlag, false, "Is client public")
	flags.Bool(internal.TrustedFlag, false, "Is the client trusted")
	flags.String(internal.DescriptionFlag, "", "Client description")
	flags.String(internal.RedirectUriFlag, "", "Redirect URIS")                       // Slice
	flags.String(internal.PostLogoutRedirectUriFlag, "", "Post logout redirect uris") // Slice

	return fctl.NewControllerConfig(
		useCreate,
		descriptionCreate,
		shortCreate,
		[]string{
			"c",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

type CreateController struct {
	store  *CreateStore
	config *fctl.ControllerConfig
}

func NewCreateController(config *fctl.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a new OAuth2 client") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	public := fctl.GetBool(flags, internal.PublicFlag)
	trusted := fctl.GetBool(flags, internal.TrustedFlag)
	description := fctl.GetString(flags, internal.DescriptionFlag)

	request := shared.CreateClientRequest{
		Public:                 &public,
		RedirectUris:           fctl.GetStringSlice(flags, internal.RedirectUriFlag),
		Description:            &description,
		Name:                   args[0],
		Trusted:                &trusted,
		PostLogoutRedirectUris: fctl.GetStringSlice(flags, internal.PostLogoutRedirectUriFlag),
	}
	response, err := authClient.Auth.CreateClient(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Client = &CreateClient{
		ID:                    response.CreateClientResponse.Data.ID,
		Name:                  response.CreateClientResponse.Data.Name,
		Description:           fctl.StringPointerToString(response.CreateClientResponse.Data.Description),
		IsPublic:              fctl.BoolPointerToString(response.CreateClientResponse.Data.Public),
		RedirectUri:           strings.Join(response.CreateClientResponse.Data.RedirectUris, ","),
		PostLogoutRedirectUri: strings.Join(response.CreateClientResponse.Data.PostLogoutRedirectUris, ","),
	}

	return c, nil
}

func (c *CreateController) Render() error {
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

func NewCreateCommand() *cobra.Command {
	config := NewCreateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](NewCreateController(config)),
	)
}
