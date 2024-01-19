package clients

import (
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// TODO: This command is a copy/paste of the create command
// We should handle membership side the patch of the client OR
// We should get the client before updating it to get replace informations

type UpdateClient struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	IsPublic              string   `json:"isPublic"`
	RedirectUri           string   `json:"redirectUri"`
	PostLogoutRedirectUri string   `json:"postLogoutRedirectUri"`
	Scopes                []string `json:"scopes"`
}

type UpdateStore struct {
	Client *UpdateClient `json:"client"`
}
type UpdateController struct {
	store                     *UpdateStore
	publicFlag                string
	trustedFlag               string
	descriptionFlag           string
	redirectUriFlag           string
	postLogoutRedirectUriFlag string
	scopes                    string
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

func NewDefaultUpdateStore() *UpdateStore {
	return &UpdateStore{
		Client: &UpdateClient{},
	}
}

func NewUpdateController() *UpdateController {
	return &UpdateController{
		store:                     NewDefaultUpdateStore(),
		publicFlag:                "public",
		trustedFlag:               "trusted",
		descriptionFlag:           "description",
		redirectUriFlag:           "redirect-uri",
		postLogoutRedirectUriFlag: "post-logout-redirect-uri",
		scopes:                    "scopes",
	}
}

func NewUpdateCommand() *cobra.Command {
	c := NewUpdateController()
	return fctl.NewCommand("update <client-id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Update client"),
		fctl.WithAliases("u", "upd"),
		fctl.WithConfirmFlag(),
		fctl.WithBoolFlag(c.publicFlag, false, "Is client public"),
		fctl.WithBoolFlag(c.trustedFlag, false, "Is the client trusted"),
		fctl.WithStringFlag(c.descriptionFlag, "", "Client description"),
		fctl.WithStringSliceFlag(c.redirectUriFlag, []string{}, "Redirect URIS"),
		fctl.WithStringSliceFlag(c.postLogoutRedirectUriFlag, []string{}, "Post logout redirect uris"),
		fctl.WithStringSliceFlag(c.scopes, []string{}, "Scopes"),
		fctl.WithController[*UpdateStore](c),
	)
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete an OAuth2 client") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	public := fctl.GetBool(cmd, c.publicFlag)
	trusted := fctl.GetBool(cmd, c.trustedFlag)
	description := fctl.GetString(cmd, c.descriptionFlag)

	request := operations.UpdateClientRequest{
		ClientID: args[0],
		UpdateClientRequest: &shared.UpdateClientRequest{
			Public:                 &public,
			RedirectUris:           fctl.GetStringSlice(cmd, c.redirectUriFlag),
			Description:            &description,
			Name:                   args[0],
			Trusted:                &trusted,
			PostLogoutRedirectUris: fctl.GetStringSlice(cmd, c.postLogoutRedirectUriFlag),
			Scopes:                 fctl.GetStringSlice(cmd, c.scopes),
		},
	}
	response, err := authClient.Auth.UpdateClient(cmd.Context(), request)
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
	c.store.Client.Scopes = response.UpdateClientResponse.Data.Scopes

	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Client.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Client.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), c.store.Client.Description})
	tableData = append(tableData, []string{pterm.LightCyan("Public"), c.store.Client.IsPublic})
	tableData = append(tableData, []string{pterm.LightCyan("Redirect URIs"), c.store.Client.RedirectUri})
	tableData = append(tableData, []string{pterm.LightCyan("Post logout redirect URIs"), c.store.Client.PostLogoutRedirectUri})
	tableData = append(tableData, []string{pterm.LightCyan("Scopes"), strings.Join(c.store.Client.Scopes, " ")})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
