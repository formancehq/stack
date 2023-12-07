package users

import (
	"fmt"
	"strconv"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateStore struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}
type UpdateController struct {
	store *UpdateStore
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

func NewDefaultUpdateStore() *UpdateStore {
	return &UpdateStore{}
}

func NewUpdateController() *UpdateController {
	return &UpdateController{
		store: NewDefaultUpdateStore(),
	}
}

func NewUpdateCommand() *cobra.Command {
	return fctl.NewCommand("update <user-id> <is-admin>",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Update user roles by by id within an organization"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*UpdateStore](NewUpdateController()),
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

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	isAdmin, err := strconv.ParseBool(args[1])
	if err != nil {
		return nil, fmt.Errorf("invalid is-admin value")
	}

	response, err := apiClient.DefaultApi.UpdateOrganizationUser(
		cmd.Context(),
		organizationID,
		args[0]).
		UpdateUserAccessData(
			membershipclient.UpdateUserAccessData{
				IsAdmin: &isAdmin,
			},
		).Execute()
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 300 {
		return nil, fmt.Errorf("error updating user: %s", response.Status)
	}
	// c.store.Id = userResponse.Data.Id
	// c.store.Email = userResponse.Data.Email
	// c.store.IsAdmin = func() bool {
	// 	if userResponse.Data.IsAdmin == nil {
	// 		return false
	// 	}

	// 	return *userResponse.Data.IsAdmin
	// }()

	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	// tableData := pterm.TableData{}
	// tableData = append(tableData, []string{pterm.LightCyan("Updated"), "true"})
	// tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.Email})

	// return pterm.DefaultTable.
	// 	WithWriter(cmd.OutOrStdout()).
	// 	WithData(tableData).
	// 	Render()
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("User updated.")
	return nil
}
