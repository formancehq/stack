package attempts

import (
	

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type RetryAllStore struct {
	ErrorResponse error `json:"error"`
	Success       bool                             `json:"success"`
}
type RetryAllController struct {
	store *RetryAllStore
}

func (c *RetryAllController) GetStore() *RetryAllStore {
	return c.store
}

var _ fctl.Controller[*RetryAllStore] = (*RetryAllController)(nil)


func (c *RetryAllController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to Retry All  Attempts") {
		return nil, fctl.ErrMissingApproval
	}

	_ , err := store.Client().Webhooks.RetryWaitingAttempts(cmd.Context())
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}
	
	return c, nil
}

func (c *RetryAllController) Render(cmd *cobra.Command, args []string) error {
	
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Attempt Retried All successfully")

	return nil
}

func NewRetryAllController() *RetryAllController {
	return &RetryAllController{
		store: &RetryAllStore{},
	}
}

func NewRetryAllCommand() *cobra.Command {
	
	return fctl.NewCommand("retry-all",
		fctl.WithShortDescription("Retry all waiting Attempts"),
		fctl.WithAliases("rtry", "rty"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*RetryAllStore](NewRetryAllController()),
	)
}
