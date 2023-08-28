package login

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useLogin         = "login"
	descriptionLogin = "Login to the service"
)

type Store struct {
	DeviceCode string `json:"deviceCode"`
	LoginURI   string `json:"loginUri"`
	BrowserURL string `json:"browserUrl"`
	Success    bool   `json:"success"`
}

func NewStore() *Store {
	return &Store{
		DeviceCode: "",
		LoginURI:   "",
		BrowserURL: "",
		Success:    false,
	}
}

func NewConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useLogin, flag.ExitOnError)
	flags.String(fctl.MembershipURIFlag, "", "service url")

	return fctl.NewControllerConfig(
		useLogin,
		descriptionLogin,
		descriptionLogin,
		[]string{
			"log",
		},
		flags,
	)
}

var _ fctl.Controller[*Store] = (*LoginController)(nil)

type LoginController struct {
	store  *Store
	config *fctl.ControllerConfig
}

func NewController(config *fctl.ControllerConfig) *LoginController {
	return &LoginController{
		store:  NewStore(),
		config: config,
	}
}

func (c *LoginController) GetStore() *Store {
	return c.store
}

func (c *LoginController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *LoginController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	membershipUri := fctl.GetString(flags, fctl.MembershipURIFlag)
	if membershipUri == "" {
		membershipUri = profile.GetMembershipURI()
	}

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(flags, map[string][]string{}, c.config.GetOut()), membershipUri)
	if err != nil {
		return nil, err
	}

	ret, err := LogIn(ctx, DialogFn(func(uri, code string) {
		c.store.DeviceCode = code
		c.store.LoginURI = uri
		fmt.Println("Link :", fmt.Sprintf("%s?user_code=%s", c.store.LoginURI, c.store.DeviceCode))
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

	currentProfileName := fctl.GetCurrentProfileName(flags, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)

	return c, cfg.Persist()
}

func (c *LoginController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Logged!")
	return nil
}

func NewCommand() *cobra.Command {
	config := NewConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*Store](NewController(config)),
	)
}
