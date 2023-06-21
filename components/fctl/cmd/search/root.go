package search

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/search/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	sizeFlag      = "size"
	defaultTarget = "ANY"
)

var targets = []string{"TRANSACTION", "ACCOUNT", "ASSET", "PAYMENT"}

type SearchStore struct {
	Response *shared.Response `json:"response"`
}

type SearchController struct {
	store  *SearchStore
	target string
}

var _ fctl.Controller[*SearchStore] = (*SearchController)(nil)

func NewDefaultSearchStore() *SearchStore {
	return &SearchStore{
		Response: &shared.Response{
			Data: make(map[string]interface{}, 0),
			Cursor: &shared.ResponseCursor{
				Data: make([]map[string]interface{}, 0),
			},
		},
	}
}

func NewSearchController() *SearchController {
	return &SearchController{
		store: NewDefaultSearchStore(),
	}
}

func (c *SearchController) GetStore() *SearchStore {
	return c.store
}

func (c *SearchController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	searchClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	terms := make([]string, 0)
	if len(args) > 1 {
		terms = args[1:]
	}
	size := int64(fctl.GetInt(cmd, sizeFlag))

	// Should be as a default value in cobra
	target := strings.ToUpper(args[0])
	if target == "" {
		target = defaultTarget
	}

	if target == "ANY" {
		target = ""
	}
	c.target = target // Save for display

	request := shared.Query{
		PageSize: &size,
		Terms:    terms,
		Target:   &target,
	}
	response, err := searchClient.Search.Search(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// "" ou ANY
	if target == "" {
		c.store.Response.Data = response.Response.Data
		c.store.Response.Cursor = response.Response.Cursor
	} else {
		// TRANSACTION, ACCOUNT, ASSET, PAYMENT
		c.store.Response.Cursor = response.Response.Cursor
	}

	return c, err
}

func (c *SearchController) Render(cmd *cobra.Command, args []string) error {
	var err error
	// No Data
	if (c.store.Response.Cursor != nil && len(c.store.Response.Cursor.Data) == 0) && len(c.store.Response.Data) == 0 {
		fctl.Section.WithWriter(cmd.OutOrStdout()).Println("No data found")
		return nil
	}

	ok := fctl.ContainValue(targets, c.target)
	// Cursor is initialized & target is valid
	if ok && c.store.Response.Cursor != nil {
		//But no data
		if len(c.store.Response.Cursor.Data) == 0 {
			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("No data found")
			return nil
		}

		// Display the data
		switch c.target {
		case "TRANSACTION":
			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Transactions")
			err = views.DisplayTransactions(cmd.OutOrStdout(), c.store.Response.Cursor.Data)
		case "ACCOUNT":
			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Accounts")
			err = views.DisplayAccounts(cmd.OutOrStdout(), c.store.Response.Cursor.Data)
		case "ASSET":
			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Assets")
			err = views.DisplayAssets(cmd.OutOrStdout(), c.store.Response.Cursor.Data)
		case "PAYMENT":
			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Payments")
			err = views.DisplayPayments(cmd.OutOrStdout(), c.store.Response.Cursor.Data)
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
			WithWriter(cmd.OutOrStdout()).
			WithData(tableData).
			Render()
	}

	return err
}

func NewCommand() *cobra.Command {

	return fctl.NewStackCommand("search <object-type> <terms>...",
		fctl.WithAliases("se"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithIntFlag(sizeFlag, 5, "Number of items to fetch"),
		fctl.WithValidArgs(append(targets, defaultTarget)...),
		// fctl.WithValidArgsFunction(func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// 	return append(targets, defaultTarget), cobra.ShellCompDirectiveNoFileComp
		// }),
		fctl.WithShortDescription("Search in all services (Default: ANY), or in a specific service (ACCOUNT, TRANSACTION, ASSET, PAYMENT)"),
		fctl.WithController[*SearchStore](NewSearchController()),
	)
}
