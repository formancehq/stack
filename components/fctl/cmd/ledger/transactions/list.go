package transactions

import (
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Transaction shared.TransactionsCursorResponseCursor `json:"transactionCursor"`
}
type ListController struct {
	store           *ListStore
	pageSizeFlag    string
	metadataFlag    string
	referenceFlag   string
	accountFlag     string
	destinationFlag string
	sourceFlag      string
	endTimeFlag     string
	startTimeFlag   string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store:           NewDefaultListStore(),
		pageSizeFlag:    "page-size",
		metadataFlag:    "metadata",
		referenceFlag:   "reference",
		accountFlag:     "account",
		destinationFlag: "dst",
		sourceFlag:      "src",
		endTimeFlag:     "end",
		startTimeFlag:   "start",
	}
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List transactions"),
		fctl.WithStringFlag(c.accountFlag, "", "Filter on account"),
		fctl.WithStringFlag(c.destinationFlag, "", "Filter on destination account"),
		fctl.WithStringFlag(c.endTimeFlag, "", "Consider transactions before date"),
		fctl.WithStringFlag(c.startTimeFlag, "", "Consider transactions after date"),
		fctl.WithStringFlag(c.sourceFlag, "", "Filter on source account"),
		fctl.WithStringFlag(c.referenceFlag, "", "Filter on reference"),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{}, "Filter transactions with metadata"),
		fctl.WithIntFlag(c.pageSizeFlag, 5, "Page size"),
		fctl.WithHiddenFlag(c.metadataFlag),
		fctl.WithArgs(cobra.ExactArgs(0)),
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

	ledger := fctl.GetString(cmd, internal.LedgerFlag)
	req := operations.ListTransactionsRequest{
		Ledger:   ledger,
		PageSize: fctl.Ptr(int64(fctl.GetInt(cmd, c.pageSizeFlag))),
	}

	if account := fctl.GetString(cmd, c.accountFlag); account != "" {
		req.Account = pointer.For(account)
	}
	if source := fctl.GetString(cmd, c.sourceFlag); source != "" {
		req.Source = pointer.For(source)
	}
	if destination := fctl.GetString(cmd, c.destinationFlag); destination != "" {
		req.Destination = pointer.For(destination)
	}
	if reference := fctl.GetString(cmd, c.referenceFlag); reference != "" {
		req.Reference = pointer.For(reference)
	}
	if startTime := fctl.GetString(cmd, c.startTimeFlag); startTime != "" {
		t, err := time.Parse(time.RFC3339Nano, startTime)
		if err != nil {
			return nil, errors.Wrap(err, "parsing start time")
		}
		req.StartTime = pointer.For(t)
	}
	if endTime := fctl.GetString(cmd, c.endTimeFlag); endTime != "" {
		t, err := time.Parse(time.RFC3339Nano, endTime)
		if err != nil {
			return nil, errors.Wrap(err, "parsing end time")
		}
		req.EndTime = pointer.For(t)
	}
	req.Metadata = collectionutils.ConvertMap(metadata, collectionutils.ToAny[string])

	response, err := ledgerClient.Ledger.ListTransactions(cmd.Context(), req)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transaction = response.TransactionsCursorResponse.Cursor

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Transaction.Data) == 0 {
		fctl.Println("No transactions found.")
		return nil
	}

	tableData := fctl.Map(c.store.Transaction.Data, func(tx shared.Transaction) []string {
		return []string{
			fmt.Sprintf("%d", tx.Txid),
			func() string {
				if tx.Reference == nil {
					return ""
				}
				return *tx.Reference
			}(),
			tx.Timestamp.Format(time.RFC3339),
			fctl.MetadataAsShortString(tx.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Reference", "Date", "Metadata"})

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
