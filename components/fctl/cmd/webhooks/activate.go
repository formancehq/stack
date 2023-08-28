package webhooks

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useActivate         = "activate <config-id>"
	descriptionActivate = "Activate one config"
	shortActivate       = "Activate one config"
)

type ActivateStore struct {
	Success bool `json:"success"`
}

func NewActivateStore() *ActivateStore {
	return &ActivateStore{
		Success: true,
	}
}
func NewActivateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useActivate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useActivate,
		descriptionActivate,
		shortActivate,
		[]string{"ac"},
		flags,
	)
}

var _ fctl.Controller[*ActivateStore] = (*Activate)(nil)

type Activate struct {
	store  *ActivateStore
	config *fctl.ControllerConfig
}

func NewActivateController(config *fctl.ControllerConfig) *Activate {
	return &Activate{
		store:  NewActivateStore(),
		config: config,
	}
}

func (c *Activate) GetStore() *ActivateStore {
	return c.store
}

func (c *Activate) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *Activate) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are bout to activate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) < 1 {
		return nil, fmt.Errorf("missing config id")
	}

	request := operations.ActivateConfigRequest{
		ID: c.config.GetArgs()[0],
	}
	response, err := client.Webhooks.ActivateConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "activating config")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return c, nil
}

func (c *Activate) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config activated successfully")

	return nil
}

func NewActivateCommand() *cobra.Command {
	config := NewActivateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ActivateStore](NewActivateController(config)),
	)
}
