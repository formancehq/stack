package instances

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InstancesSendEventStore struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instance_id"`
	Event      string `json:"event"`
}
type InstancesSendEventController struct {
	store *InstancesSendEventStore
}

var _ fctl.Controller[*InstancesSendEventStore] = (*InstancesSendEventController)(nil)

func NewDefaultInstancesSendEventStore() *InstancesSendEventStore {
	return &InstancesSendEventStore{
		Success:    false,
		InstanceID: "",
		Event:      "",
	}
}

func NewInstancesSendEventController() *InstancesSendEventController {
	return &InstancesSendEventController{
		store: NewDefaultInstancesSendEventStore(),
	}
}

func NewSendEventCommand() *cobra.Command {
	return fctl.NewCommand("send-event <instance-id> <event>",
		fctl.WithShortDescription("Send an event to an instance"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*InstancesSendEventStore](NewInstancesSendEventController()),
	)
}

func (c *InstancesSendEventController) GetStore() *InstancesSendEventStore {
	return c.store
}

func (c *InstancesSendEventController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}
	response, err := client.Orchestration.SendEvent(cmd.Context(), operations.SendEventRequest{
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

func (c *InstancesSendEventController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Event '%s' sent", args[1])
	return nil
}
