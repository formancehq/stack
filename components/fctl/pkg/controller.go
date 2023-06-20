package fctl

import (
	"github.com/spf13/cobra"
)

type Renderable interface {
	Render(cmd *cobra.Command, args []string) error
}
type Controller interface {
	GetStore() *SharedStore
	Run(cmd *cobra.Command, args []string) (Renderable, error)
}

func WithRender(cmd *cobra.Command, args []string, c Controller, r Renderable) error {
	flags := GetString(cmd, OutputFlag)

	switch flags {
	case "json":
		data, err := c.GetStore().ToJson()
		if (err) != nil {
			return err
		}
		cmd.OutOrStdout().Write(data)
		return nil
	case "plain":
		return r.Render(cmd, args)
	}

	return nil
}
