package cmd

import (
	"context"
	"fmt"
	"net/url"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
)

type Dialog interface {
	DisplayURIAndCode(uri, code string)
}
type DialogFn func(uri, code string)

func (fn DialogFn) DisplayURIAndCode(uri, code string) {
	fn(uri, code)
}

func LogIn(ctx context.Context, dialog Dialog, relyingParty rp.RelyingParty) (*oidc.Tokens, error) {
	deviceCode, err := rp.GetDeviceCode(ctx, relyingParty)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(deviceCode.GetVerificationUri())
	if err != nil {
		panic(err)
	}
	query := uri.Query()
	query.Set("user_code", deviceCode.GetUserCode())
	uri.RawQuery = query.Encode()

	dialog.DisplayURIAndCode(deviceCode.GetVerificationUri(), deviceCode.GetUserCode())

	if err := fctl.Open(uri.String()); err != nil {
		return nil, err
	}

	return rp.PollDeviceCode(ctx, deviceCode.GetDeviceCode(), deviceCode.GetInterval(), relyingParty)
}

type LoginOutput struct {
	profile    *fctl.Profile `json:"-"`
	DeviceCode string        `json:"device_code"`
	LoginURI   string        `json:"login_uri"`
	BrowserURL string        `json:"browser_url"`
	Success    bool          `json:"success"`
}

func loginCommand(cmd *cobra.Command, args []string) error {

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

	relyingParty, err := fctl.GetAuthRelyingParty(cmd, membershipUri)
	if err != nil {
		return err
	}

	loginOutput := &LoginOutput{
		profile: profile,
	}

	ret, err := LogIn(cmd.Context(), DialogFn(func(uri, code string) {
		loginOutput.DeviceCode = code
		loginOutput.LoginURI = uri
	}), relyingParty)

	// Other relying error not related to browser
	if err != nil && err.Error() != "error_opening_browser" {
		return err
	}

	// Browser not found
	if err != nil && err.Error() == "error_opening_browser" {
		url, ok := fctl.GetSharedAdditionnalData("browser_url").(string)
		if !ok {
			return err
		}

		loginOutput.BrowserURL = url
	} else {
		loginOutput.Success = true
	}

	profile.SetMembershipURI(membershipUri)
	profile.UpdateToken(ret.Token)

	currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)
	fctl.SetSharedData(loginOutput, profile, cfg, nil)

	return cfg.Persist()
}

func displayLoginResult(cmd *cobra.Command, args []string) error {

	data := fctl.GetSharedData().(*LoginOutput)

	fmt.Println("Please enter the following code on your browser:", data.DeviceCode)
	fmt.Println("Link:", data.LoginURI)

	if !data.Success && data.BrowserURL != "" {
		fmt.Printf("Unable to find a browser, please open the following link: %s", data.BrowserURL)
		return nil
	}

	if data.Success {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Logged!")
	}

	return nil

}

func NewLoginCommand() *cobra.Command {
	return fctl.NewCommand("login",
		fctl.WithStringFlag(fctl.MembershipURIFlag, "", "service url"),
		fctl.WithHiddenFlag(fctl.MembershipURIFlag),
		fctl.WithShortDescription("Login"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(loginCommand),
		fctl.WrapOutputPostRunE(displayLoginResult),
	)
}
