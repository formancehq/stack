package clients

import (
	"flag"
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List clients"
)

type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
	IsPublic    string `json:"isPublic"`
}

type ListStore struct {
	Clients []Client `json:"clients"`
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
		fctl.Organization,
		fctl.Stack,
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

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	clients, err := authClient.Auth.ListClients(ctx)
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

func (c *ListController) Render() error {
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
