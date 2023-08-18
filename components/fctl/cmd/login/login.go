package login

import (
	"errors"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Dialog interface {
	DisplayURIAndCode(err error, uri, code string)
}
type DialogFn func(err error, uri, code string)

func (fn DialogFn) DisplayURIAndCode(err error, uri, code string) {
	fn(err, uri, code)
}

type LoginStore struct {
	profile    *fctl.Profile `json:"-"`
	DeviceCode string        `json:"deviceCode"`
	LoginURI   string        `json:"loginUri"`
	BrowserURL string        `json:"browserUrl"`
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

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(cmd, map[string][]string{}), membershipUri)
	if err != nil {
		return nil, err
	}

	c.store.profile = profile

	ret, err := LogIn(cmd.Context(), DialogFn(func(err error, uri, code string) {
		c.store.DeviceCode = code
		c.store.LoginURI = uri
		if errors.Is(err, fctl.ErrOpenningBrowser) {
			fmt.Println("No browser detected")
			fmt.Println("Link :", fmt.Sprintf("%s?user_code=%s", c.store.LoginURI, c.store.DeviceCode))
		}
	}), relyingParty)

	// Other relying error not related to browser
	if err != nil {
		return nil, err
	}

	// Browser not found
	if ret != nil {
		c.store.Success = true
		profile.UpdateToken(ret)
	}

	profile.SetMembershipURI(membershipUri)

	currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)

	return c, cfg.Persist()
}

func (c *LoginController) Render(cmd *cobra.Command, args []string) error {
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
