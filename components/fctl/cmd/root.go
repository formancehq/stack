package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/formancehq/fctl/cmd/auth"
	"github.com/formancehq/fctl/cmd/cloud"
	"github.com/formancehq/fctl/cmd/ledger"
	"github.com/formancehq/fctl/cmd/login"
	"github.com/formancehq/fctl/cmd/orchestration"
	"github.com/formancehq/fctl/cmd/payments"
	"github.com/formancehq/fctl/cmd/profiles"
	"github.com/formancehq/fctl/cmd/search"
	"github.com/formancehq/fctl/cmd/stack"
	"github.com/formancehq/fctl/cmd/ui"
	"github.com/formancehq/fctl/cmd/version"
	"github.com/formancehq/fctl/cmd/wallets"
	"github.com/formancehq/fctl/cmd/webhooks"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cmd := fctl.NewCommand("fctl",
		fctl.WithSilenceError(),
		fctl.WithShortDescription("Formance Control CLI"),
		fctl.WithSilenceUsage(),
		fctl.WithChildCommands(
			ui.NewCommand(),
			version.NewCommand(),
			login.NewCommand(),
			NewPromptCommand(),
			ledger.NewCommand(),
			payments.NewCommand(),
			profiles.NewCommand(),
			stack.NewCommand(),
			auth.NewCommand(),
			cloud.NewCommand(),
			search.NewCommand(),
			webhooks.NewCommand(),
			wallets.NewCommand(),
			orchestration.NewCommand(),
		),
		fctl.WithPersistentStringPFlag(fctl.ProfileFlag, "p", "", "config profile to use"),
		fctl.WithPersistentStringPFlag(fctl.FileFlag, "c", fmt.Sprintf("%s/.formance/fctl.config", homedir), "Debug mode"),
		fctl.WithPersistentBoolPFlag(fctl.DebugFlag, "d", false, "Debug mode"),
		fctl.WithPersistentStringPFlag(fctl.OutputFlag, "o", "plain", "Output format (plain, json)"),
		fctl.WithPersistentBoolFlag(fctl.InsecureTlsFlag, false, "Insecure TLS"),
		fctl.WithPersistentBoolFlag(fctl.TelemetryFlag, false, "Telemetry enabled"),
	)
	// cmd.Version = Version
	return cmd
}

func Execute() {
	defer func() {
		if e := recover(); e != nil {
			pterm.Error.WithWriter(os.Stderr).Printfln("%s", e)
			debug.PrintStack()
		}
	}()

	ctx, _ := signal.NotifyContext(context.TODO(), os.Interrupt)
	err := NewRootCommand().ExecuteContext(ctx)
	if err != nil {
		switch {
		case errors.Is(err, fctl.ErrMissingApproval):
			pterm.Error.WithWriter(os.Stderr).Printfln("Command aborted as you didn't approve.")
			os.Exit(1)
		case extractOpenAPIErrorMessage(err) != nil:
			pterm.Error.WithWriter(os.Stderr).Printfln(extractOpenAPIErrorMessage(err).Error())
			os.Exit(2)
		default:
			pterm.Error.WithWriter(os.Stderr).Printfln(err.Error())
			os.Exit(255)
		}
	}
}

func extractOpenAPIErrorMessage(err error) error {
	if err == nil {
		return nil
	}

	if err := unwrapOpenAPIError(err); err != nil {
		return errors.New(err.ErrorMessage)
	}

	return err
}

func unwrapOpenAPIError(err error) *shared.ErrorResponse {
	openapiError := &membershipclient.GenericOpenAPIError{}
	if errors.As(err, &openapiError) {
		body := openapiError.Body()
		// Actually, each api redefine errors response
		// So OpenAPI generator generate an error structure for every service
		// Manually unmarshal errorResponse allow us to handle only one ErrorResponse
		// It will be refined once the monorepo fully ready
		errResponse := api.ErrorResponse{}
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil
		}

		if errResponse.ErrorCode != "" {
			errorCode := shared.ErrorsEnum(errResponse.ErrorCode)
			return &shared.ErrorResponse{
				ErrorCode:    errorCode,
				ErrorMessage: errResponse.ErrorMessage,
				Details:      &errResponse.Details,
			}
		}
	}

	return nil
}
