package accounts

import (
	"fmt"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Accounts []shared.Account `json:"accounts"`
}
type ListController struct {
	store        *ListStore
	metadataFlag string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store:        NewDefaultListStore(),
		metadataFlag: "metadata",
	}
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List accounts"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{}, "Filter accounts with metadata"),
		fctl.WithController[*ListStore](c),
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

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	request := operations.ListAccountsRequest{
		Ledger:   fctl.GetString(cmd, internal.LedgerFlag),
		Metadata: metadata,
	}
	rsp, err := ledgerClient.Ledger.ListAccounts(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if rsp.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", rsp.ErrorResponse.ErrorCode, rsp.ErrorResponse.ErrorMessage)
	}

	if rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	c.store.Accounts = rsp.AccountsCursorResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {

	tableData := fctl.Map(c.store.Accounts, func(account shared.Account) []string {
		return []string{
			account.Address,
			fctl.MetadataAsShortString(account.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Address", "Metadata"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
