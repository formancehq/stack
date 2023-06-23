package billing

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PortalStore struct {
	FoundBrowser   bool   `json:"found_browser"`
	BillingPlanUrl string `json:"billing_plan_url"`
}
type PortalController struct {
	store *PortalStore
}

var _ fctl.Controller[*PortalStore] = (*PortalController)(nil)

func NewDefaultPortalStore() *PortalStore {
	return &PortalStore{
		BillingPlanUrl: "",
		FoundBrowser:   true,
	}
}

func NewPortalController() *PortalController {
	return &PortalController{
		store: NewDefaultPortalStore(),
	}
}

func NewPortalCommand() *cobra.Command {
	return fctl.NewCommand("portal",
		fctl.WithAliases("p"),
		fctl.WithShortDescription("Access to Billing Portal"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*PortalStore](NewPortalController()),
	)
}

func (c *PortalController) GetStore() *PortalStore {
	return c.store
}

func (c *PortalController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	billing, _, err := apiClient.DefaultApi.BillingPortal(cmd.Context(), organizationID).Execute()
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

func (c *PortalController) Render(cmd *cobra.Command, args []string) error {
	if c.store.FoundBrowser && c.store.BillingPlanUrl != "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Billing Portal opened in your browser")
		return nil
	}

	if !c.store.FoundBrowser && c.store.BillingPlanUrl != "" {
		pterm.Error.WithWriter(cmd.OutOrStdout()).Printfln("Please open %s in your browser", c.store.BillingPlanUrl)
		return nil
	}

	if c.store.BillingPlanUrl == "" {
		pterm.Error.WithWriter(cmd.OutOrStdout()).Printfln("Please subscribe to a plan to access Billing Portal")
		return nil
	}

	return nil
}
