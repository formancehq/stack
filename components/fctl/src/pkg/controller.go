package fctl

import (
	"github.com/spf13/cobra"
)

type Renderable interface {
	Render(cmd *cobra.Command, args []string) error
}
type Controller[T any] interface {
	GetStore() T
	Run(cmd *cobra.Command, args []string) (Renderable, error)
}
type ExportedData struct {
	Data interface{} `json:"data"`
}
