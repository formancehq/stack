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

type LoginOutput struct {
	profile    *fctl.Profile `json:"-"`
	DeviceCode string        `json:"device_code"`
	LoginURI   string        `json:"login_uri"`
	BrowserURL string        `json:"browser_url"`
	Success    bool          `json:"success"`
}
type LoginController struct {
	store *fctl.SharedStore
}

func (c *LoginController) GetStore() *fctl.SharedStore {
	return c.store
}
func NewLoginController() *LoginController {
	return &LoginController{
		store: fctl.NewSharedStore(),
	}
}
func (c *LoginController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	membershipUri, err := cmd.Flags().GetString(fctl.MembershipURIFlag)
	if err != nil {
		return nil, err
	}
	if membershipUri == "" {
		membershipUri = profile.GetMembershipURI()
	}

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(cmd, map[string][]string{}), membershipUri)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Browser not found
	if err != nil && err.Error() == "error_opening_browser" {
		// url, ok := fctl.GetSharedAdditionnalData("browser_url").(string)
		// if !ok {
		// 	return err
		// }

		// loginOutput.BrowserURL = url
	} else {
		loginOutput.Success = true
	}

	profile.SetMembershipURI(membershipUri)
	profile.UpdateToken(ret)

	currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)
	// fctl.SetSharedData(loginOutput, profile, cfg, nil)

	return c, cfg.Persist()
}

func (c *LoginController) Render(cmd *cobra.Command, args []string) error {

	data := c.store.GetData().(*LoginOutput)

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
		fctl.WithController(NewLoginController()),
	)
}
