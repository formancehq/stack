package attempts

import (
	

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type AbortStore struct {
	ErrorResponse error `json:"error"`
	Success       bool                             `json:"success"`
}
type AbortController struct {
	store *AbortStore
}

func (c *AbortController) GetStore() *AbortStore {
	return c.store
}

var _ fctl.Controller[*AbortStore] = (*AbortController)(nil)


func (c *AbortController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	request := operations.AbortWaitingAttemptRequest{
		AttemptID: args[0],
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to abort an Attempt") {
		return nil, fctl.ErrMissingApproval
	}

	_ , err := store.Client().Webhooks.AbortWaitingAttempt(cmd.Context(), request)
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}
	
	return c, nil
}

func (c *AbortController) Render(cmd *cobra.Command, args []string) error {
	
	
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}


	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Attempt abort successfully")

	return nil
}

func NewAbortController() *AbortController {
	return &AbortController{
		store: &AbortStore{},
	}
}

func NewAbortCommand() *cobra.Command {
	
	return fctl.NewCommand("abort <attempt-id>",
		fctl.WithShortDescription("Abort a waiting Attempt"),
		fctl.WithAliases("abrt", "ab"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*AbortStore](NewAbortController()),
	)
}
