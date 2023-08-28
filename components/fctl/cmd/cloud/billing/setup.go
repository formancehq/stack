package billing

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useSetup   = "setup"
	shortSetup = "Create a new billing account"
)

type SetupStore struct {
	FoundBrowser bool   `json:"foundBrowser"`
	BillingUrl   string `json:"billingUrl"`
}

func NewSetupStore() *SetupStore {
	return &SetupStore{
		BillingUrl:   "",
		FoundBrowser: true,
	}
}

func NewSetupConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSetup, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useSetup,
		shortSetup,
		shortSetup,
		[]string{
			"setup", "set",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*SetupStore] = (*SetupController)(nil)

type SetupController struct {
	store  *SetupStore
	config *fctl.ControllerConfig
}

func NewSetupController(config *fctl.ControllerConfig) *SetupController {
	return &SetupController{
		store:  NewSetupStore(),
		config: config,
	}
}

func (c *SetupController) GetStore() *SetupStore {
	return c.store
}

func (c *SetupController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *SetupController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	billing, _, err := apiClient.DefaultApi.BillingSetup(ctx, organizationID).Execute()
	if err != nil {
		pterm.Error.WithWriter(c.config.GetOut()).Printfln("You already have an active subscription")
		return c, nil
	}

	c.store.BillingUrl = billing.Data.Url
	if err := fctl.Open(billing.Data.Url); err != nil {
		c.store.FoundBrowser = false
	}

	return c, nil
}

func (c *SetupController) Render() error {

	out := c.config.GetOut()

	if !c.store.FoundBrowser && c.store.BillingUrl != "" {
		pterm.Warning.WithWriter(out).Printfln("Could not open browser, please visit %s", c.store.BillingUrl)
		return nil
	}

	if c.store.BillingUrl == "" {
		pterm.Error.WithWriter(out).Printfln("You already have an active subscription")
		return nil
	}

	if c.store.FoundBrowser && c.store.BillingUrl != "" {
		pterm.Success.WithWriter(out).Printfln("Billing Setup opened in your browser")
		return nil
	}

	return nil
}
func NewSetupCommand() *cobra.Command {
	config := NewSetupConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithDeprecated("Please contact Formances Sales Team."),
		fctl.WithController[*SetupStore](NewSetupController(config)),
	)
}
