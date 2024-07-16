package ledger

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ImportStore struct{}
type ImportController struct {
	store         *ImportStore
	inputFileFlag string
}

var _ fctl.Controller[*ImportStore] = (*ImportController)(nil)

func NewDefaultImportStore() *ImportStore {
	return &ImportStore{}
}

func NewImportController() *ImportController {
	return &ImportController{
		store:         NewDefaultImportStore(),
		inputFileFlag: "file",
	}
}

func NewImportCommand() *cobra.Command {
	c := NewImportController()
	return fctl.NewCommand("import <ledger name> <file path>",
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithShortDescription("Import a ledger"),
		fctl.WithStringFlag(c.inputFileFlag, "", "Import from stdin or file"),
		fctl.WithController[*ImportStore](c),
	)
}

func (c *ImportController) GetStore() *ImportStore {
	return c.store
}

func (c *ImportController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	_, err := store.Client().Ledger.V2ImportLogs(cmd.Context(), operations.V2ImportLogsRequest{
		Ledger:      args[0],
		RequestBody: pointer.For(fmt.Sprintf("file:%s", args[1])),
	})

	return c, err
}

func (c *ImportController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Ledger imported!")
	return nil
}
