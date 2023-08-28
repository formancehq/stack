package instances

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
	useSendEvent   = "send-event <instance-id> <event>"
	shortSendEvent = "Send an event to an instance"
)

type SendEventStore struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instanceId"`
	Event      string `json:"event"`
}

func NewSendEventConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSendEvent, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useSendEvent,
		shortSendEvent,
		shortSendEvent,
		[]string{
			"se",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}

type SendEventController struct {
	store  *SendEventStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*SendEventStore] = (*SendEventController)(nil)

func NewSendEventStore() *SendEventStore {
	return &SendEventStore{
		Success:    false,
		InstanceID: "",
		Event:      "",
	}
}

func NewSendEventController(config *fctl.ControllerConfig) *SendEventController {
	return &SendEventController{
		store:  NewSendEventStore(),
		config: config,
	}
}

func (c *SendEventController) GetStore() *SendEventStore {
	return c.store
}

func (c *SendEventController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *SendEventController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfig(flags, ctx, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}
	response, err := client.Orchestration.SendEvent(ctx, operations.SendEventRequest{
		RequestBody: &operations.SendEventRequestBody{
			Name: args[1],
		},
		InstanceID: args[0],
	})

	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.InstanceID = args[0]
	c.store.Event = args[1]

	return c, nil
}

func (c *SendEventController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Event '%s' sent", c.config.GetArgs()[1])
	return nil
}

func NewSendEventCommand() *cobra.Command {
	config := NewSendEventConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*SendEventStore](NewSendEventController(config)),
	)
}
