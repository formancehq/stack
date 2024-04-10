package volumes

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/printer"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Cursor shared.V2VolumesWithBalanceCursorResponseCursor
}

type ListController struct {
	store                *ListStore
	cursorFlag           string
	pageSizeFlag         string
	metadataFlag         string
	addressFlag          string
	pitFlag              string
	ootFlag              string
	useInsertionDateFlag string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	address := fctl.GetString(cmd, c.addressFlag)

	body := make([]map[string]map[string]any, 0)

	for key, value := range metadata {
		body = append(body, map[string]map[string]any{
			"$match": {
				"metadata[" + key + "]": value,
			},
		})
	}

	if address != "" {
		body = append(body, map[string]map[string]any{
			"$match": {
				"account": address,
			},
		})
	}

	cursor := fctl.GetString(cmd, c.cursorFlag)
	pageSize := int64(fctl.GetInt(cmd, c.pageSizeFlag))
	useInsertionDate := fctl.GetBool(cmd, c.useInsertionDateFlag)

	pit, err := fctl.GetDateTime(cmd, c.pitFlag)
	if err != nil {
		return nil, err
	}
	oot, err := fctl.GetDateTime(cmd, c.ootFlag)
	if err != nil {
		return nil, err
	}

	if pit == nil {
		tmp := time.Now()
		pit = &tmp
	}

	if oot == nil {
		tmp := pit.Add(-24 * 7 * time.Hour)
		oot = &tmp
	}

	request := operations.V2GetVolumesWithBalancesRequest{
		RequestBody:   collectionutils.ConvertMap(metadata, collectionutils.ToAny[string]),
		Ledger:        fctl.GetString(cmd, internal.LedgerFlag),
		StartTime:     oot,
		EndTime:       pit,
		Cursor:        &cursor,
		PageSize:      &pageSize,
		InsertionDate: &useInsertionDate,
	}

	response, err := store.Client().Ledger.V2GetVolumesWithBalances(cmd.Context(), request)

	if err != nil {
		return nil, errors.Wrap(err, "Get Volumes With Balances")
	}

	c.store.Cursor = response.V2VolumesWithBalanceCursorResponse.Cursor

	return c, nil

}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {

	tableData := fctl.Map(c.store.Cursor.Data, func(volume shared.V2VolumesWithBalance) []string {
		return []string{
			volume.Account,
			volume.Asset,
			fmt.Sprintf("%d", volume.Input),
			fmt.Sprintf("%d", volume.Output),
			fmt.Sprintf("%d", volume.Balance),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"Account", "Asset", "Input", "Output", "Balance"})

	tableData = printer.AddCursorRowsToTable(tableData, printer.CursorArgs{
		HasMore:  c.store.Cursor.HasMore,
		Next:     c.store.Cursor.Next,
		PageSize: c.store.Cursor.PageSize,
		Previous: c.store.Cursor.Previous,
	})

	writer := cmd.OutOrStdout()

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(writer).
		WithData(tableData).
		Render()

}

func NewListController() *ListController {
	return &ListController{
		store:                &ListStore{},
		metadataFlag:         "metadata",
		cursorFlag:           "cursor",
		pageSizeFlag:         "page-size",
		addressFlag:          "address",
		pitFlag:              "end-time",
		ootFlag:              "start-time",
		useInsertionDateFlag: "insertion-date",
	}
}

func NewListCommand() *cobra.Command {

	c := NewListController()

	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List volumes and balances for a period of time (OOT-PIT)"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithStringFlag(c.pitFlag, "", "PIT (Point in Time)"),
		fctl.WithStringFlag(c.ootFlag, "", "OOT (Origin of Time)"),
		fctl.WithBoolFlag(c.useInsertionDateFlag, false, "Use insertion date"),
		fctl.WithStringFlag(c.addressFlag, "", "Filter accounts with address"),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{}, "Filter accounts with metadata"),
		fctl.WithStringFlag(c.cursorFlag, "", "Cursor pagination"),
		fctl.WithIntFlag(c.pageSizeFlag, 10, "Page size"),
		fctl.WithController[*ListStore](c),
	)
}
