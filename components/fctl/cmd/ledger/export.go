package ledger

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/spf13/cobra"
)

type ExportStore struct {
	response *http.Response
}
type ExportController struct {
	store          *ExportStore
	outputFileFlag string
}

var _ fctl.Controller[*ExportStore] = (*ExportController)(nil)

func NewDefaultExportStore() *ExportStore {
	return &ExportStore{}
}

func NewExportController() *ExportController {
	return &ExportController{
		store:          NewDefaultExportStore(),
		outputFileFlag: "file",
	}
}

func NewExportCommand() *cobra.Command {
	c := NewExportController()
	return fctl.NewCommand("export",
		fctl.WithShortDescription("Export a ledger"),
		fctl.WithStringFlag(c.outputFileFlag, "", "Export to file"),
		fctl.WithController[*ExportStore](c),
	)
}

func (c *ExportController) GetStore() *ExportStore {
	return c.store
}

func (c *ExportController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	ctx := cmd.Context()
	out := fctl.GetString(cmd, "file")
	if out != "" {
		ctx = context.WithValue(ctx, "path", out)
	}

	ret, err := store.Client().Ledger.V2.ExportLogs(ctx, operations.V2ExportLogsRequest{
		Ledger: fctl.GetString(cmd, internal.LedgerFlag),
	})
	if err != nil {
		return nil, err
	}

	c.store.response = ret.RawResponse

	return c, nil
}

func (c *ExportController) Render(cmd *cobra.Command, args []string) error {
	out := fctl.GetString(cmd, "file")
	if out == "" {
		_, err := io.Copy(os.Stdout, c.store.response.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
