package hooks

import (
	"net/url"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	secretFlag = "secret"
	nameFlag = "name"
	retryFlag = "retry"
)

type CreateWebhookController struct {
	store *CreateWebhookStore
}

type CreateWebhookStore struct {
	Webhook shared.V2Hook `json:"webhook"`
	ErrorResponse error `json:"error"`
	
}

var _ fctl.Controller[*CreateWebhookStore] = (*CreateWebhookController)(nil)

func NewDefaultCreateWebhookStore() *CreateWebhookStore {
	return &CreateWebhookStore{
		Webhook: shared.V2Hook{},
	}
}

func NewCreateWebhookController() *CreateWebhookController {
	return &CreateWebhookController{
		store: NewDefaultCreateWebhookStore(),
	}
}
func (c *CreateWebhookController) GetStore() *CreateWebhookStore {
	return c.store
}

func (c *CreateWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to create a webhook") {
		return nil, fctl.ErrMissingApproval
	}
	
	Endpoint :=  args[0]
	Events := args[1:]

	if _, err := url.Parse(Endpoint); err != nil {
		return nil, errors.Wrap(err, "invalid endpoint URL")
	}

	secret := fctl.GetString(cmd, secretFlag)
	name := fctl.GetString(cmd, nameFlag)
	retry := fctl.GetBool(cmd, retryFlag)

	hookBodyParams := shared.V2HookBodyParams{
		Endpoint: Endpoint,
		Events : Events,
		Secret: &secret,
		Retry: &retry,
		Name : &name,
	}

	response, err := store.Client().Webhooks.InsertHook(cmd.Context(), hookBodyParams)
	
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Webhook = response.V2HookResponse.Data
	}

	return c, nil
}

func (c *CreateWebhookController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hook ID : %s created successfully", c.store.Webhook.ID)
	return nil
}

func NewCreateCommand() *cobra.Command {

	return fctl.NewCommand("create <endpoint> [<event-type>...]",
		fctl.WithShortDescription("Create a new Hook. At least one event type is required."),
		fctl.WithAliases("cr"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithStringFlag(secretFlag, "", "Bring your own webhooks signing secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)"),
		fctl.WithStringFlag(nameFlag, "", "Name for the Hook (optionnal)"),
		fctl.WithBoolFlag(retryFlag, true, "Does the hook should retry failed attempt (Default:true)"),		
		fctl.WithController[*CreateWebhookStore](NewCreateWebhookController()),
	)
}
