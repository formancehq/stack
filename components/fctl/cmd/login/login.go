package login

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useLogin         = "login"
	descriptionLogin = "Login to the service"
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
	DeviceCode string        `json:"deviceCode"`
	LoginURI   string        `json:"loginUri"`
	BrowserURL string        `json:"browserUrl"`
	Success    bool          `json:"success"`
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

type LoginControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewLoginControllerConfig() *LoginControllerConfig {
	flags := flag.NewFlagSet(useLogin, flag.ExitOnError)
	flags.String(fctl.MembershipURIFlag, "", "service url")
	fctl.WithGlobalFlags(flags)

	return &LoginControllerConfig{
		context:     nil,
		use:         useLogin,
		description: descriptionLogin,
		aliases:     []string{},
		out:         os.Stdout,
		flags:       flags,
		args:        []string{},
	}
}

var _ fctl.Controller[*LoginStore] = (*LoginController)(nil)

type LoginController struct {
	store  *LoginStore
	config LoginControllerConfig
}

func NewLoginController(config LoginControllerConfig) *LoginController {
	return &LoginController{
		store:  NewDefaultLoginStore(),
		config: config,
	}
}

func (c *LoginController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *LoginController) GetContext() context.Context {
	return c.config.context
}

func (c *LoginController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *LoginController) GetStore() *LoginStore {
	return c.store
}

func (c *LoginController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *LoginController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	membershipUri := fctl.GetString(flags, fctl.MembershipURIFlag)
	if membershipUri == "" {
		membershipUri = profile.GetMembershipURI()
	}

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(flags, map[string][]string{}), membershipUri)
	if err != nil {
		return nil, err
	}

	c.store.profile = profile

	ret, err := LogIn(ctx, DialogFn(func(uri, code string) {
		c.store.DeviceCode = code
		c.store.LoginURI = uri
	}), relyingParty)

	// Other relying error not related to browser
	if err != nil && err.Error() != "error_opening_browser" {
		return nil, err
	}

	// Browser not found
	if err == nil {
		c.store.Success = true
	}

	profile.SetMembershipURI(membershipUri)
	profile.UpdateToken(ret)

	currentProfileName := fctl.GetCurrentProfileName(flags, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)

	return c, cfg.Persist()
}

func (c *LoginController) Render() error {

	fmt.Println("Please enter the following code on your browser:", c.store.DeviceCode)
	fmt.Println("Link:", c.store.LoginURI)

	if !c.store.Success && c.store.BrowserURL != "" {
		fmt.Printf("Unable to find a browser, please open the following link: %s", c.store.BrowserURL)
		return nil
	}

	if c.store.Success {
		pterm.Success.WithWriter(c.config.out).Printfln("Logged!")
	}

	return nil

}

func NewCommand() *cobra.Command {
	config := NewLoginControllerConfig()
	return fctl.NewCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*LoginStore](NewLoginController(*config)),
	)
}
