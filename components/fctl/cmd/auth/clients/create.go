package clients

import (
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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
type CreateController struct {
	store                     *CreateStore
	publicFlag                string
	trustedFlag               string
	descriptionFlag           string
	redirectUriFlag           string
	postLogoutRedirectUriFlag string
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store:                     NewDefaultCreateStore(),
		publicFlag:                "public",
		trustedFlag:               "trusted",
		descriptionFlag:           "description",
		redirectUriFlag:           "redirect-uri",
		postLogoutRedirectUriFlag: "post-logout-redirect-uri",
	}
}

func NewCreateCommand() *cobra.Command {
	c := NewCreateController()
	return fctl.NewCommand("create <name>",
		fctl.WithAliases("c"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithBoolFlag(c.publicFlag, false, "Is client public"),
		fctl.WithBoolFlag(c.trustedFlag, false, "Is the client trusted"),
		fctl.WithStringFlag(c.descriptionFlag, "", "Client description"),
		fctl.WithStringSliceFlag(c.redirectUriFlag, []string{}, "Redirect URIS"),
		fctl.WithStringSliceFlag(c.postLogoutRedirectUriFlag, []string{}, "Post logout redirect uris"),
		fctl.WithShortDescription("Create client"),
		fctl.WithController[*CreateStore](c),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to create a new OAuth2 client") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	public := fctl.GetBool(cmd, c.publicFlag)
	trusted := fctl.GetBool(cmd, c.trustedFlag)
	description := fctl.GetString(cmd, c.descriptionFlag)

	request := shared.CreateClientRequest{
		Public:                 &public,
		RedirectUris:           fctl.GetStringSlice(cmd, c.redirectUriFlag),
		Description:            &description,
		Name:                   args[0],
		Trusted:                &trusted,
		PostLogoutRedirectUris: fctl.GetStringSlice(cmd, c.postLogoutRedirectUriFlag),
	}
	response, err := authClient.Auth.CreateClient(cmd.Context(), request)
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

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Client.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Client.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), c.store.Client.Description})
	tableData = append(tableData, []string{pterm.LightCyan("Public"), c.store.Client.IsPublic})
	tableData = append(tableData, []string{pterm.LightCyan("Redirect URIs"), c.store.Client.RedirectUri})
	tableData = append(tableData, []string{pterm.LightCyan("Post logout redirect URIs"), c.store.Client.PostLogoutRedirectUri})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
