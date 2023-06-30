package install

import (
	"fmt"
	"io/ioutil"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewBankingCircleCommand() *cobra.Command {
	const (
		endpointFlag                   = "endpoint"
		authorizationEndpointFlag      = "authorization-endpoint"
		userCertificateFilePathFlag    = "user-certificate"
		userCertificateKeyFilePathFlag = "user-certificate-key"

		defaultEndpoint              = "https://sandbox.bankingcircle.com"
		defaultAuthorizationEndpoint = "https://authorizationsandbox.bankingcircleconnect.com"
	)
	return fctl.NewCommand(internal.BankingCircleConnector+" <username> <password>",
		fctl.WithShortDescription("Install a Banking Circle connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(endpointFlag, defaultEndpoint, "API endpoint"),
		fctl.WithStringFlag(authorizationEndpointFlag, defaultAuthorizationEndpoint, "Authorization endpoint"),
		fctl.WithStringFlag(userCertificateFilePathFlag, "", "User certificate"),
		fctl.WithStringFlag(userCertificateKeyFilePathFlag, "", "User certificate key"),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to install connector '%s'", internal.BankingCircleConnector) {
				return fctl.ErrMissingApproval
			}

			paymentsClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			certificate, err := ioutil.ReadFile(fctl.GetString(cmd, userCertificateFilePathFlag))
			if err != nil {
				return errors.Wrap(err, "reading user certificate")
			}

			key, err := ioutil.ReadFile(fctl.GetString(cmd, userCertificateKeyFilePathFlag))
			if err != nil {
				return errors.Wrap(err, "reading user certificate key")
			}

			request := operations.InstallConnectorRequest{
				Connector: shared.ConnectorBankingCircle,
				RequestBody: shared.BankingCircleConfig{
					Username:              args[0],
					Password:              args[1],
					Endpoint:              fctl.GetString(cmd, endpointFlag),
					AuthorizationEndpoint: fctl.GetString(cmd, authorizationEndpointFlag),
					UserCertificate:       string(certificate),
					UserCertificateKey:    string(key),
				},
			}
			response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "installing connector")
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector installed!")

			return nil
		}),
	)
}
