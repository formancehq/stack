package policies

import (
	"encoding/json"
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Policy shared.Policy `json:"policy"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("get <policyID>",
		fctl.WithShortDescription("Get policy"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := ledgerClient.Reconciliation.GetPolicy(cmd.Context(), operations.GetPolicyRequest{
		PolicyID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if response.PolicyResponse == nil {
		return nil, fmt.Errorf("policy not found")
	}

	c.store.Policy = response.PolicyResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Policy.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Policy.Name})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.Policy.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("LedgerName"), c.store.Policy.LedgerName})
	tableData = append(tableData, []string{pterm.LightCyan("LedgerQuey"), func() string {
		if c.store.Policy.LedgerQuery == nil {
			return ""
		}
		raw, _ := json.Marshal(c.store.Policy.LedgerQuery)
		return string(raw)
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("PaymentsPoolID"), c.store.Policy.PaymentsPoolID})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}
