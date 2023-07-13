package fctl

import (
	"context"
	"flag"
)

type Renderable interface {
	Render() error
}
type Controller[T any] interface {
	GetStore() T
	GetFlags() *flag.FlagSet
	SetArgs([]string)

	// Not present from creation of the controller
	// The context is set by the command
	GetContext() context.Context
	SetContext(context.Context)

	Run() (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}
