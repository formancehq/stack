package fctl

import (
	"flag"
)

type Renderable interface {
	Render() error
}
type Controller[T any] interface {
	GetStore() T
	GetFlags() *flag.FlagSet
	Run() (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}
