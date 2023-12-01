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

type PaymentsConnectorsCurrencyCloudStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsCurrencyCloudController struct {
	store                *PaymentsConnectorsCurrencyCloudStore
	endpointFlag         string
	defaultEndpoint      string
	pollingPeriodFlag    string
	defaultpollingPeriod string
	nameFlag             string
	defaultName          string
}

var _ fctl.Controller[*PaymentsConnectorsCurrencyCloudStore] = (*PaymentsConnectorsCurrencyCloudController)(nil)

func NewDefaultPaymentsConnectorsCurrencyCloudStore() *PaymentsConnectorsCurrencyCloudStore {
	return &PaymentsConnectorsCurrencyCloudStore{
		Success: false,
	}
}

func NewPaymentsConnectorsCurrencyCloudController() *PaymentsConnectorsCurrencyCloudController {
	return &PaymentsConnectorsCurrencyCloudController{
		store:                NewDefaultPaymentsConnectorsCurrencyCloudStore(),
		endpointFlag:         "endpoint",
		defaultEndpoint:      "https://devapi.currencycloud.com",
		pollingPeriodFlag:    "polling-period",
		defaultpollingPeriod: "2m",
		nameFlag:             "name",
		defaultName:          "currencycloud",
	}
}

func NewCurrencyCloudCommand() *cobra.Command {
	c := NewPaymentsConnectorsCurrencyCloudController()
	return fctl.NewCommand(internal.CurrencyCloudConnector+" <login-id> <api-key>",
		fctl.WithShortDescription("Install a Currency Cloud connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithStringFlag(c.pollingPeriodFlag, c.defaultpollingPeriod, "Polling duration"),
		fctl.WithStringFlag(c.nameFlag, c.defaultName, "Connector name"),
		fctl.WithController[*PaymentsConnectorsCurrencyCloudStore](c),
	)
}

func (c *PaymentsConnectorsCurrencyCloudController) GetStore() *PaymentsConnectorsCurrencyCloudStore {
	return c.store
}

func (c *PaymentsConnectorsCurrencyCloudController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.CurrencyCloudConnector)
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
		RequestBody: &shared.CurrencyCloudConfig{
			Name:          fctl.GetString(cmd, c.nameFlag),
			APIKey:        args[1],
			LoginID:       args[0],
			Endpoint:      endpoint,
			PollingPeriod: fctl.Ptr(fctl.GetString(cmd, c.pollingPeriodFlag)),
		},
		Connector: shared.ConnectorCurrencyCloud,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.CurrencyCloudConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsCurrencyCloudController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
