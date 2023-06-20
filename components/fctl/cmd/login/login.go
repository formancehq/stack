package login

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Dialog interface {
	DisplayURIAndCode(uri, code string)
}
type DialogFn func(uri, code string)

func (fn DialogFn) DisplayURIAndCode(uri, code string) {
	fn(uri, code)
}

type LoginStore struct {
	profile    *fctl.Profile `json:"-"`
	DeviceCode string        `json:"device_code"`
	LoginURI   string        `json:"login_uri"`
	BrowserURL string        `json:"browser_url"`
	Success    bool          `json:"success"`
}
type LoginController struct {
	store *LoginStore
}

func NewDefaultLoginStore() *LoginStore {
	return &LoginStore{
		profile:    nil,
		DeviceCode: "",
		LoginURI:   "",
		BrowserURL: "",
		Success:    false,
	}
}
func (c *LoginController) GetStore() *LoginStore {
	return c.store
}
func NewLoginController() *LoginController {
	return &LoginController{
		store: NewDefaultLoginStore(),
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

	relyingParty, err := fctl.GetAuthRelyingParty(cmd, membershipUri)
	if err != nil {
		return nil, err
	}

	c.store.profile = profile

	ret, err := LogIn(cmd.Context(), DialogFn(func(uri, code string) {
		c.store.DeviceCode = code
		c.store.LoginURI = uri
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
		c.store.Success = true
	}

	profile.SetMembershipURI(membershipUri)
	profile.UpdateToken(ret.Token)

	currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)
	// fctl.SetSharedData(loginOutput, profile, cfg, nil)

	return c, cfg.Persist()
}

func (c *LoginController) Render(cmd *cobra.Command, args []string) error {

	fmt.Println("Please enter the following code on your browser:", c.store.DeviceCode)
	fmt.Println("Link:", c.store.LoginURI)

	if !c.store.Success && c.store.BrowserURL != "" {
		fmt.Printf("Unable to find a browser, please open the following link: %s", c.store.BrowserURL)
		return nil
	}

	if c.store.Success {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Logged!")
	}

	return nil

}

func NewCommand() *cobra.Command {
	return fctl.NewCommand("login",
		fctl.WithStringFlag(fctl.MembershipURIFlag, "", "service url"),
		fctl.WithHiddenFlag(fctl.MembershipURIFlag),
		fctl.WithShortDescription("Login"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*LoginStore](NewLoginController()),
	)
}
