package attempts

import (
	

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type RetryStore struct {
	ErrorResponse error `json:"error"`
	Success       bool  `json:"success"`
}
type RetryController struct {
	store *RetryStore
}

func (c *RetryController) GetStore() *RetryStore {
	return c.store
}

var _ fctl.Controller[*RetryStore] = (*RetryController)(nil)


func (c *RetryController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	request := operations.RetryWaitingAttemptRequest{
		AttemptID: args[0],
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to Retry an Attempt") {
		return nil, fctl.ErrMissingApproval
	}

	_ , err := store.Client().Webhooks.RetryWaitingAttempt(cmd.Context(), request)
	
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}
	
	return c, nil
}

func (c *RetryController) Render(cmd *cobra.Command, args []string) error {
	
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Attempt Retry successfully")

	return nil
}

func NewRetryController() *RetryController {
	return &RetryController{
		store: &RetryStore{},
	}
}

func NewRetryCommand() *cobra.Command {
	
	return fctl.NewCommand("retry <attempt-id>",
		fctl.WithShortDescription("Retry a waiting Attempt"),
		fctl.WithAliases("rtry", "rty"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*RetryStore](NewRetryController()),
	)
}
