package holds

import (
	"flag"
	"fmt"
	"math/big"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useConfirm   = "confirm <hold-id>"
	shortConfirm = "Confirm a hold"
	finalFlag    = "final"
	amountFlag   = "amount"
)

type ConfirmStore struct {
	Success bool   `json:"success"`
	HoldId  string `json:"holdId"`
}
type ConfirmController struct {
	store  *ConfirmStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ConfirmStore] = (*ConfirmController)(nil)

func NewConfirmStore() *ConfirmStore {
	return &ConfirmStore{}
}
func NewConfirmConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useConfirm, flag.ExitOnError)
	flags.Bool(finalFlag, false, "Is final debit (close hold)")
	flags.Int(amountFlag, 0, "Amount to confirm")

	c := fctl.NewControllerConfig(
		useConfirm,
		shortConfirm,
		shortConfirm,
		[]string{
			"c", "conf",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}
func NewConfirmController(config *fctl.ControllerConfig) *ConfirmController {
	return &ConfirmController{
		store:  NewConfirmStore(),
		config: config,
	}
}

func (c *ConfirmController) GetStore() *ConfirmStore {
	return c.store
}

func (c *ConfirmController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ConfirmController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	final := fctl.GetBool(flags, finalFlag)
	amount := int64(fctl.GetInt(flags, amountFlag))

	request := operations.ConfirmHoldRequest{
		HoldID: c.config.GetArgs()[0],
		ConfirmHoldRequest: &shared.ConfirmHoldRequest{
			Amount: big.NewInt(amount),
			Final:  &final,
		},
	}
	response, err := stackClient.Wallets.ConfirmHold(ctx, request)
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
	c.store.HoldId = c.config.GetArgs()[0]

	return c, nil
}

func (c *ConfirmController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Hold '%s' confirmed!", c.config.GetArgs()[0])
	return nil
}

func NewConfirmCommand() *cobra.Command {
	c := NewConfirmConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithController[*ConfirmStore](NewConfirmController(c)),
	)
}
