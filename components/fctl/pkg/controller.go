package fctl

import (
	"context"
	"flag"
	"io"
)

type Renderable interface {
	Render() error
}

type Controller[T any] interface {
	GetStore() T

	GetConfig() *ControllerConfig

	Run() (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}

type ControllerConfig struct {
	context          context.Context
	use              string
	description      string
	shortDescription string
	aliases          []string
	out              io.Writer
	flags            *flag.FlagSet
	pflags           *flag.FlagSet
	scopes           *flag.FlagSet
	args             []string
}

func NewControllerConfig(use string, description string, shortDescription string, aliases []string, flagSet *flag.FlagSet, scopes ...*flag.Flag) *ControllerConfig {
	return &ControllerConfig{
		use:              use,
		description:      description,
		shortDescription: shortDescription,
		aliases:          aliases,
		flags:            flagSet,
		scopes:           WithScopesFlags(flag.NewFlagSet("scopes", flag.ExitOnError), scopes...),
		pflags:           GlobalFlags,
	}

}

func (c *ControllerConfig) GetScopes() *flag.FlagSet {
	return c.scopes
}

func (c *ControllerConfig) GetUse() string {
	return c.use
}

func (c *ControllerConfig) GetDescription() string {
	return c.description
}

func (c *ControllerConfig) GetShortDescription() string {
	return c.shortDescription
}

func (c *ControllerConfig) SetShortDescription(shortDescription string) {
	c.shortDescription = shortDescription
}

func (c *ControllerConfig) GetAliases() []string {
	return c.aliases
}

func (c *ControllerConfig) GetOut() io.Writer {
	return c.out
}
func (c *ControllerConfig) SetOut(out io.Writer) {
	c.out = out
}

func (c *ControllerConfig) GetArgs() []string {
	return c.args
}

func (c *ControllerConfig) SetArgs(args []string) {
	c.args = args
}

func (c *ControllerConfig) GetFlags() *flag.FlagSet {
	return c.flags
}

// GetAllFLags Return the pflags & flags merged together in a new FlagSet
// This is done to avoid mutating the original flag.FlagSet
// which is used by the controller to parse the flags
// and the pflags are used by the controller to parse the persistent one
func (c *ControllerConfig) GetAllFLags() *flag.FlagSet {

	// Create a new FlagSet
	flags := flag.NewFlagSet(c.use, flag.ExitOnError)

	// Regroup pflag // GLOBAL
	if c.pflags != nil {
		c.pflags.VisitAll(func(f *flag.Flag) {
			flags.Var(f.Value, f.Name, f.Usage)
		})
	}

	// Regroup flags
	if c.flags != nil {
		c.flags.VisitAll(func(f *flag.Flag) {
			flags.Var(f.Value, f.Name, f.Usage)
		})
	}

	// Regroup scopes
	if c.scopes != nil {
		c.scopes.VisitAll(func(f *flag.Flag) {
			flags.Var(f.Value, f.Name, f.Usage)
		})
	}

	return flags
}

func (c *ControllerConfig) GetPFlags() *flag.FlagSet {
	return c.pflags
}

func (c *ControllerConfig) GetContext() context.Context {
	if c.context == nil {
		return context.TODO()
	}

	return c.context
}

func (c *ControllerConfig) SetContext(ctx context.Context) {
	c.context = ctx
}
