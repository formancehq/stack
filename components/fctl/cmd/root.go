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
	"github.com/formancehq/fctl/cmd/reconciliation"
	"github.com/formancehq/fctl/cmd/search"
	"github.com/formancehq/fctl/cmd/stack"
	"github.com/formancehq/fctl/cmd/ui"
	"github.com/formancehq/fctl/cmd/version"
	"github.com/formancehq/fctl/cmd/wallets"
	"github.com/formancehq/fctl/cmd/webhooks"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
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
			reconciliation.NewCommand(),
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
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			logger := logging.NewLogrus(logging.DefaultLogger(cmd.OutOrStdout(), fctl.GetBool(cmd, fctl.DebugFlag)))
			ctx := logging.ContextWithLogger(cmd.Context(), logger)
			cmd.SetContext(ctx)
			return nil
		}),
	)

	cmd.Version = version.Version
	cmd.RegisterFlagCompletionFunc(fctl.ProfileFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := fctl.GetConfig(cmd)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		ret := make([]string, 0)
		for name := range cfg.GetProfiles() {
			ret = append(ret, name)
		}
		return ret, cobra.ShellCompDirectiveDefault
	})
	return cmd
}

func Execute() {
	defer func() {
		if e := recover(); e != nil {
			pterm.Error.WithWriter(os.Stderr).Printfln("%s", e)
			debug.PrintStack()
		}
	}()
	cobra.EnableTraverseRunHooks = true
	ctx, _ := signal.NotifyContext(context.TODO(), os.Interrupt)
	err := NewRootCommand().ExecuteContext(ctx)
	if err != nil {
		switch {
		case errors.Is(err, fctl.ErrMissingApproval):
			pterm.Error.WithWriter(os.Stderr).Printfln("Command aborted as you didn't approve.")
			os.Exit(1)
		case fctl.IsInvalidAuthentication(err):
			pterm.Error.WithWriter(os.Stderr).Printfln("Your authentication is invalid, please login :)")
		default:
			unwrapped := errors.Unwrap(err)
			for unwrapped != nil {
				//notes(gfyrag): not a clean assertion but following errors does not implements standard Is() helper for errors
				switch err := unwrapped.(type) {
				case *sdkerrors.ErrorResponse:
					printErrorResponse(err)
					return
				case *sdkerrors.V2ErrorResponse:
					printV2ErrorResponse(err)
					return
				case *membershipclient.GenericOpenAPIError:
					body := err.Body()
					errResponse := api.ErrorResponse{}
					if err := json.Unmarshal(body, &errResponse); err != nil {
						panic(err)
					}
					printError(errResponse.ErrorCode, errResponse.ErrorMessage, &errResponse.Details)
					return
				default:
					unwrapped = errors.Unwrap(unwrapped)
				}
			}
		}
		pterm.Error.WithWriter(os.Stderr).Printfln(err.Error())
		os.Exit(255)
	}
}

func printError(code string, message string, details *string) {
	pterm.Error.WithWriter(os.Stderr).Printfln("Got error with code %s: %s", code, message)
	if details != nil && *details != "" {
		pterm.Error.WithWriter(os.Stderr).Printfln("Details:\r\n%s", *details)
	}
	os.Exit(2)
}

func printV2ErrorResponse(target *sdkerrors.V2ErrorResponse) {
	printError(string(target.ErrorCode), target.ErrorMessage, target.Details)
}

func printErrorResponse(target *sdkerrors.ErrorResponse) {
	printError(string(target.ErrorCode), target.ErrorMessage, target.Details)
}
