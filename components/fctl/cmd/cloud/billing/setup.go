package billing

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SetupStore struct {
	FoundBrowser bool   `json:"found_browser"`
	BillingUrl   string `json:"billing_url"`
}
type SetupController struct {
	store *SetupStore
}

var _ fctl.Controller[*SetupStore] = (*SetupController)(nil)

func NewDefaultSetupStore() *SetupStore {
	return &SetupStore{
		BillingUrl:   "",
		FoundBrowser: true,
	}
}

func NewSetupController() *SetupController {
	return &SetupController{
		store: NewDefaultSetupStore(),
	}
}

func NewSetupCommand() *cobra.Command {
	return fctl.NewCommand("setup",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Create a new billing account"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*SetupStore](NewSetupController()),
	)
}

func (c *SetupController) GetStore() *SetupStore {
	return c.store
}

func (c *SetupController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	billing, _, err := apiClient.DefaultApi.BillingSetup(cmd.Context(), organizationID).Execute()
	if err != nil {
		pterm.Error.WithWriter(cmd.OutOrStderr()).Printfln("You already have an active subscription")
		return c, nil
	}

	c.store.BillingUrl = billing.Data.Url
	if err := fctl.Open(billing.Data.Url); err != nil {
		c.store.FoundBrowser = false
	}

	return c, nil
}

func (c *SetupController) Render(cmd *cobra.Command, args []string) error {
	if !c.store.FoundBrowser && c.store.BillingUrl != "" {
		pterm.Warning.WithWriter(cmd.OutOrStderr()).Printfln("Could not open browser, please visit %s", c.store.BillingUrl)
		return nil
	}

	if c.store.BillingUrl == "" {
		pterm.Error.WithWriter(cmd.OutOrStderr()).Printfln("You already have an active subscription")
		return nil
	}

	if c.store.FoundBrowser && c.store.BillingUrl != "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Billing Setup opened in your browser")
		return nil
	}

	return nil
}
