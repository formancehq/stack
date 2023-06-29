package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

type Dialog interface {
	DisplayURIAndCode(uri, code string)
}
type DialogFn func(uri, code string)

func (fn DialogFn) DisplayURIAndCode(uri, code string) {
	fn(uri, code)
}

func LogIn(ctx context.Context, dialog Dialog, relyingParty rp.RelyingParty) (*oidc.AccessTokenResponse, error) {
	deviceCode, err := rp.DeviceAuthorization(relyingParty.OAuthConfig().Scopes, relyingParty)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(deviceCode.VerificationURI)
	if err != nil {
		panic(err)
	}
	query := uri.Query()
	query.Set("user_code", deviceCode.UserCode)
	uri.RawQuery = query.Encode()

	dialog.DisplayURIAndCode(deviceCode.VerificationURI, deviceCode.UserCode)

	if err := fctl.Open(uri.String()); err != nil {
		return nil, err
	}

	return rp.DeviceAccessToken(ctx, deviceCode.DeviceCode, time.Duration(deviceCode.Interval)*time.Second, relyingParty)
}

func NewLoginCommand() *cobra.Command {
	return fctl.NewCommand("login",
		fctl.WithStringFlag(fctl.MembershipURIFlag, "", "service url"),
		fctl.WithHiddenFlag(fctl.MembershipURIFlag),
		fctl.WithShortDescription("Login"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {

			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			profile := fctl.GetCurrentProfile(cmd, cfg)
			membershipUri, err := cmd.Flags().GetString(fctl.MembershipURIFlag)
			if err != nil {
				return err
			}
			if membershipUri == "" {
				membershipUri = profile.GetMembershipURI()
			}

			relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(cmd, map[string][]string{}), membershipUri)
			if err != nil {
				return err
			}

			ret, err := LogIn(cmd.Context(), DialogFn(func(uri, code string) {
				fmt.Fprintln(cmd.OutOrStdout(), "Please enter the following code on your browser:", code)
				fmt.Fprintln(cmd.OutOrStdout(), "Link:", uri)
			}), relyingParty)
			if err != nil {
				return err
			}

			profile.SetMembershipURI(membershipUri)
			profile.UpdateToken(ret)

			currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

			cfg.SetCurrentProfile(currentProfileName, profile)

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Logged!")
			return cfg.Persist()
		}),
	)
}
