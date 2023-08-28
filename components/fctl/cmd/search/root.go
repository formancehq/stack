package search

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/search/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useSearch         = "search <object-type> <terms>..."
	shortDescription  = "Search for transactions, accounts, assets, and payments"
	descriptionSearch = "Search in all services (Default: ANY), or in a specific service (ACCOUNT, TRANSACTION, ASSET, PAYMENT)"
	sizeFlag          = "size"
	defaultTarget     = "ANY"
)

var targets = []string{"TRANSACTION", "ACCOUNT", "ASSET", "PAYMENT"}

type Store struct {
	Response *shared.Response `json:"response"`
}

func NewSearchConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSearch, flag.ExitOnError)
	flags.Int(sizeFlag, 5, "Number of items to fetch")
	return fctl.NewControllerConfig(
		useSearch,
		descriptionSearch,
		shortDescription,
		[]string{
			"se",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

func NewStore() *Store {
	return &Store{
		Response: &shared.Response{
			Data: make(map[string]interface{}, 0),
			Cursor: &shared.ResponseCursor{
				Data: make([]map[string]interface{}, 0),
			},
		},
	}
}

var _ fctl.Controller[*Store] = (*Controller)(nil)

type Controller struct {
	store  *Store
	target string
	config *fctl.ControllerConfig
}

func NewController(config *fctl.ControllerConfig) *Controller {
	return &Controller{
		store:  NewStore(),
		config: config,
	}
}

func (c *Controller) GetStore() *Store {
	return c.store
}

func (c *Controller) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *Controller) Run() (fctl.Renderable, error) {

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

	searchClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	terms := make([]string, 0)
	if len(args) > 1 {
		terms = args[1:]
	}
	size := int64(fctl.GetInt(flags, sizeFlag))

	target := strings.ToUpper(args[0])

	if target == "ANY" {
		target = ""
	}
	c.target = target
	request := shared.Query{
		PageSize: &size,
		Terms:    terms,
		Target:   &target,
	}
	response, err := searchClient.Search.Search(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if target == "" {
		c.store.Response.Data = response.Response.Data
		c.store.Response.Cursor = response.Response.Cursor
	} else {
		// TRANSACTION, ACCOUNT, ASSET, PAYMENT
		c.store.Response.Cursor = response.Response.Cursor
	}

	return c, err
}

func (c *Controller) Render() error {
	var err error
	out := c.config.GetOut()
	// No Data
	if (c.store.Response.Cursor != nil && len(c.store.Response.Cursor.Data) == 0) && len(c.store.Response.Data) == 0 {
		fctl.Section.WithWriter(out).Println("No data found")
		return nil
	}

	ok := fctl.ContainValue(targets, c.target)
	// Cursor is initialized & target is valid
	if ok && c.store.Response.Cursor != nil {
		//But no data
		if len(c.store.Response.Cursor.Data) == 0 {
			fctl.Section.WithWriter(out).Println("No data found")
			return nil
		}

		// Display the data
		switch c.target {
		case "TRANSACTION":
			fctl.Section.WithWriter(out).Println("Transactions")
			err = views.DisplayTransactions(out, c.store.Response.Cursor.Data)
		case "ACCOUNT":
			fctl.Section.WithWriter(out).Println("Accounts")
			err = views.DisplayAccounts(out, c.store.Response.Cursor.Data)
		case "ASSET":
			fctl.Section.WithWriter(out).Println("Assets")
			err = views.DisplayAssets(out, c.store.Response.Cursor.Data)
		case "PAYMENT":
			fctl.Section.WithWriter(out).Println("Payments")
			err = views.DisplayPayments(out, c.store.Response.Cursor.Data)
		}
	}

	ok = defaultTarget == c.target || c.target == ""

	// Any data
	if len(c.store.Response.Data) > 0 && ok {
		tableData := make([][]string, 0)
		for kind, values := range c.store.Response.Data {
			for _, value := range values.([]any) {
				dataAsJson, err := json.Marshal(value)
				if err != nil {
					return err
				}

				dataAsJsonString := string(dataAsJson)
				if len(dataAsJsonString) > 100 {
					dataAsJsonString = dataAsJsonString[:100] + "..."
				}

				tableData = append(tableData, []string{
					kind, dataAsJsonString,
				})
			}
		}
		tableData = fctl.Prepend(tableData, []string{"Kind", "Object"})
		return pterm.DefaultTable.
			WithHasHeader().
			WithWriter(out).
			WithData(tableData).
			Render()
	}

	return err
}

func NewCommand() *cobra.Command {
	config := NewSearchConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithValidArgs(append(targets, defaultTarget)...),
		fctl.WithController[*Store](NewController(config)),
	)
}
