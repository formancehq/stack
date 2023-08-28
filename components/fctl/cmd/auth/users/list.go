package users

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List all users subjects"
)

type User struct {
	ID      string `json:"id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Email   string `json:"email,omitempty"`
}

type ListStore struct {
	Users []User `json:"users,omitempty"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}

func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
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

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	listUsersResponse, err := client.Auth.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	if listUsersResponse.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", listUsersResponse.StatusCode)
	}

	c.store.Users = fctl.Map(listUsersResponse.ListUsersResponse.Data, func(o shared.User) User {
		return User{
			ID:      *o.ID,
			Subject: *o.Subject,
			Email:   *o.Email,
		}
	})

	return c, nil
}

func (c *ListController) Render() error {
	if len(c.store.Users) == 0 {
		fmt.Fprintln(c.config.GetOut(), "No users found.")
		return nil
	}

	tableData := fctl.Map(c.store.Users, func(o User) []string {
		return []string{
			o.ID,
			o.Subject,
			o.Email,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Subject", "Email"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
