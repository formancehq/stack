package install

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsConnectorsModulrStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}

type PaymentsConnectorsModulrController struct {
	store                *PaymentsConnectorsModulrStore
	endpointFlag         string
	defaultEndpoint      string
	pollingPeriodFlag    string
	defaultpollingPeriod string
	nameFlag             string
	defaultName          string
}

var _ fctl.Controller[*PaymentsConnectorsModulrStore] = (*PaymentsConnectorsModulrController)(nil)

func NewDefaultPaymentsConnectorsModulrStore() *PaymentsConnectorsModulrStore {
	return &PaymentsConnectorsModulrStore{
		Success: false,
	}
}

func NewPaymentsConnectorsModulrController() *PaymentsConnectorsModulrController {
	return &PaymentsConnectorsModulrController{
		store:                NewDefaultPaymentsConnectorsModulrStore(),
		endpointFlag:         "endpoint",
		defaultEndpoint:      "https://api-sandbox.modulrfinance.com",
		pollingPeriodFlag:    "polling-period",
		defaultpollingPeriod: "2m",
		nameFlag:             "name",
		defaultName:          "modulr",
	}
}

func NewModulrCommand() *cobra.Command {
	c := NewPaymentsConnectorsModulrController()
	return fctl.NewCommand(internal.ModulrConnector+" <api-key> <api-secret>",
		fctl.WithShortDescription("Install a Modulr connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithStringFlag(c.pollingPeriodFlag, c.defaultpollingPeriod, "Polling duration"),
		fctl.WithStringFlag(c.nameFlag, c.defaultName, "Connector name"),
		fctl.WithController[*PaymentsConnectorsModulrStore](c),
	)
}

func (c *PaymentsConnectorsModulrController) GetStore() *PaymentsConnectorsModulrStore {
	return c.store
}

func (c *PaymentsConnectorsModulrController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.ModulrConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	var endpoint *string
	if e := fctl.GetString(cmd, c.endpointFlag); e != "" {
		endpoint = &e
	}

	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
		RequestBody: &shared.ModulrConfig{
			Name:          fctl.GetString(cmd, c.nameFlag),
			APIKey:        args[0],
			APISecret:     args[1],
			Endpoint:      endpoint,
			PollingPeriod: fctl.Ptr(fctl.GetString(cmd, c.pollingPeriodFlag)),
		},
		Connector: shared.ConnectorModulr,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.ModulrConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsModulrController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
