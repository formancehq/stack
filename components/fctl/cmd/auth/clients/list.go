package clients

import (
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
	IsPublic    string `json:"is_public"`
}

type ListStore struct {
	Clients []Client `json:"clients"`
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
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List clients"),
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

	authClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	clients, err := authClient.Auth.ListClients(cmd.Context())
	if err != nil {
		return nil, err
	}

	if clients.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", clients.StatusCode)
	}

	c.store.Clients = fctl.Map(clients.ListClientsResponse.Data, func(o shared.Client) Client {
		return Client{
			ID:   o.ID,
			Name: o.Name,
			Description: func() string {
				if o.Description == nil {
					return ""
				}
				return ""
			}(),
			Scope:    strings.Join(o.Scopes, ","),
			IsPublic: fctl.BoolPointerToString(o.Public),
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Clients, func(o Client) []string {
		return []string{
			o.ID,
			o.Name,
			o.Description,
			o.Scope,
			o.IsPublic,
		}
	})

	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Description", "Scopes", "Public"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
