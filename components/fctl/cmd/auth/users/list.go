package users

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type User struct {
	ID      string `json:"id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Email   string `json:"email,omitempty"`
}

type ListStore struct {
	Users []User `json:"users,omitempty"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List users"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	listUsersResponse, err := client.Auth.ListUsers(cmd.Context())
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

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Users) == 0 {
		fctl.Println("No users found.")
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
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
