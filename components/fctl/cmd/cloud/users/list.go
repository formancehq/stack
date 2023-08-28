package users

import (
	"flag"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List users"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type ListStore struct {
	Users []User `json:"users"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"l", "ls",
		},
		flags,
		fctl.Organization,
	)
}

func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	usersResponse, _, err := apiClient.DefaultApi.ListUsers(ctx, organizationID).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Users = fctl.Map(usersResponse.Data, func(i membershipclient.User) User {
		return User{
			i.Id,
			i.Email,
		}
	})

	return c, nil
}

func (c *ListController) Render() error {

	usersRow := fctl.Map(c.store.Users, func(i User) []string {
		return []string{
			i.ID,
			i.Email,
		}
	})

	tableData := fctl.Prepend(usersRow, []string{"ID", "Email"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}
func NewListCommand() *cobra.Command {

	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
