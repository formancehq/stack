package billing

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	usePortal   = "portal"
	shortPortal = "Open the billing portal"
)

type PortalStore struct {
	FoundBrowser   bool   `json:"foundBrowser"`
	BillingPlanUrl string `json:"billingPlanUrl"`
}

func NewPortalStore() *PortalStore {
	return &PortalStore{
		BillingPlanUrl: "",
		FoundBrowser:   true,
	}
}

func NewPortalConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(usePortal, flag.ExitOnError)
	return fctl.NewControllerConfig(
		usePortal,
		shortPortal,
		shortPortal,
		[]string{
			"p",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*PortalStore] = (*PortalController)(nil)

type PortalController struct {
	store  *PortalStore
	config *fctl.ControllerConfig
}

func NewPortalController(config *fctl.ControllerConfig) *PortalController {
	return &PortalController{
		store:  NewPortalStore(),
		config: config,
	}
}

func (c *PortalController) GetStore() *PortalStore {
	return c.store
}

func (c *PortalController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *PortalController) Run() (fctl.Renderable, error) {

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

	billing, _, err := apiClient.DefaultApi.BillingPortal(ctx, organizationID).Execute()
	if err != nil {
		return nil, err
	}
	if billing != nil {
		c.store.BillingPlanUrl = billing.Data.Url
		if err := fctl.Open(billing.Data.Url); err != nil {
			c.store.FoundBrowser = false
		}
	}

	return c, nil
}

func (c *PortalController) Render() error {
	out := c.config.GetOut()

	if c.store.FoundBrowser && c.store.BillingPlanUrl != "" {
		pterm.Success.WithWriter(out).Printfln("Billing Portal opened in your browser")
		return nil
	}

	if !c.store.FoundBrowser && c.store.BillingPlanUrl != "" {
		pterm.Error.WithWriter(out).Printfln("Please open %s in your browser", c.store.BillingPlanUrl)
		return nil
	}

	if c.store.BillingPlanUrl == "" {
		pterm.Error.WithWriter(out).Printfln("Please subscribe to a plan to access Billing Portal")
		return nil
	}

	return nil
}

func NewPortalCommand() *cobra.Command {
	config := NewPortalConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithDeprecated("Please contact Formances Sales Team."),
		fctl.WithController[*PortalStore](NewPortalController(config)),
	)
}
