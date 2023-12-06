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

type PaymentsConnectorsMoneycorpStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsMoneycorpController struct {
	store                *PaymentsConnectorsMoneycorpStore
	endpointFlag         string
	defaultEndpoint      string
	pollingPeriodFlag    string
	defaultpollingPeriod string
	nameFlag             string
	defaultName          string
}

func NewDefaultPaymentsConnectorsMoneycorpStore() *PaymentsConnectorsMoneycorpStore {
	return &PaymentsConnectorsMoneycorpStore{
		Success:       false,
		ConnectorName: internal.MoneycorpConnector,
	}
}
func NewPaymentsConnectorsMoneycorpController() *PaymentsConnectorsMoneycorpController {
	return &PaymentsConnectorsMoneycorpController{
		store:                NewDefaultPaymentsConnectorsMoneycorpStore(),
		endpointFlag:         "endpoint",
		defaultEndpoint:      "https://sandbox-corpapi.moneycorp.com",
		pollingPeriodFlag:    "polling-period",
		defaultpollingPeriod: "2m",
		nameFlag:             "name",
		defaultName:          "moneycorp",
	}
}
func NewMoneycorpCommand() *cobra.Command {
	c := NewPaymentsConnectorsMoneycorpController()

	return fctl.NewCommand(internal.MoneycorpConnector+" <clientID> <apiKey>",
		fctl.WithShortDescription("Install a Moneycorp connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithStringFlag(c.pollingPeriodFlag, c.defaultpollingPeriod, "Polling duration"),
		fctl.WithStringFlag(c.nameFlag, c.defaultName, "Connector name"),
		fctl.WithController[*PaymentsConnectorsMoneycorpStore](c),
	)
}

func (c *PaymentsConnectorsMoneycorpController) GetStore() *PaymentsConnectorsMoneycorpStore {
	return c.store
}

func (c *PaymentsConnectorsMoneycorpController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to install connector '%s'", internal.MoneycorpConnector) {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMoneycorp,
		ConnectorConfig: shared.ConnectorConfig{
			MoneycorpConfig: &shared.MoneycorpConfig{
				Name:          fctl.GetString(cmd, c.nameFlag),
				ClientID:      args[0],
				APIKey:        args[1],
				Endpoint:      fctl.GetString(cmd, c.endpointFlag),
				PollingPeriod: fctl.Ptr(fctl.GetString(cmd, c.pollingPeriodFlag)),
			},
		},
	}
	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector installed!")

	c.store.Success = true
	c.store.ConnectorName = internal.MoneycorpConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsMoneycorpController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
