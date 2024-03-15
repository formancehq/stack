package holds

import (
	"fmt"
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ConfirmStore struct {
	Success bool   `json:"success"`
	HoldId  string `json:"holdId"`
}
type ConfirmController struct {
	store      *ConfirmStore
	finalFlag  string
	amountFlag string
}

var _ fctl.Controller[*ConfirmStore] = (*ConfirmController)(nil)

func NewDefaultConfirmStore() *ConfirmStore {
	return &ConfirmStore{}
}

func NewConfirmController() *ConfirmController {
	return &ConfirmController{
		store:      NewDefaultConfirmStore(),
		finalFlag:  "final",
		amountFlag: "amount",
	}
}

func NewConfirmCommand() *cobra.Command {
	c := NewConfirmController()
	return fctl.NewCommand("confirm <hold-id>",
		fctl.WithShortDescription("Confirm a hold"),
		fctl.WithAliases("c", "conf"),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithBoolFlag(c.finalFlag, false, "Is final debit (close hold)"),
		fctl.WithIntFlag(c.amountFlag, 0, "Amount to confirm"),
		fctl.WithController[*ConfirmStore](c),
	)
}

func (c *ConfirmController) GetStore() *ConfirmStore {
	return c.store
}

func (c *ConfirmController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())

	final := fctl.GetBool(cmd, c.finalFlag)
	amount := int64(fctl.GetInt(cmd, c.amountFlag))

	request := operations.ConfirmHoldRequest{
		HoldID: args[0],
		ConfirmHoldRequest: &shared.ConfirmHoldRequest{
			Amount: big.NewInt(amount),
			Final:  &final,
		},
	}
	response, err := store.Client().Wallets.ConfirmHold(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "confirming hold")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true //Todo: check status code
	c.store.HoldId = args[0]

	return c, nil
}

func (c *ConfirmController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hold '%s' confirmed!", args[0])

	return nil

}
