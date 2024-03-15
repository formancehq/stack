package ledger

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	bucketNameFlag = "bucket"
)

type CreateStore struct{}

type CreateController struct {
	store         *CreateStore
	metadataFlag  string
	referenceFlag string
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store:         NewDefaultCreateStore(),
		metadataFlag:  "metadata",
		referenceFlag: "reference",
	}
}

func NewCreateCommand() *cobra.Command {
	c := NewCreateController()
	return fctl.NewCommand("create <name>",
		fctl.WithAliases("c", "cr"),
		fctl.WithShortDescription("Create a new ledger (starting from ledger v2)"),
		fctl.WithStringFlag(bucketNameFlag, "", "Bucket on which install the new ledger"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](c),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to create a new ledger") {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Ledger.V2CreateLedger(cmd.Context(), operations.V2CreateLedgerRequest{
		V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
			Bucket: pointer.For(fctl.GetString(cmd, bucketNameFlag)),
		},
		Ledger: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.V2ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.V2ErrorResponse.ErrorCode, response.V2ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d when creating ledger", response.StatusCode)
	}

	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Ledger created!")
	return nil
}
