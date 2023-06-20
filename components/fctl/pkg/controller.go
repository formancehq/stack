package fctl

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
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

func WithRender[T any](cmd *cobra.Command, args []string, c Controller[T], r Renderable) error {
	flags := GetString(cmd, OutputFlag)

	switch flags {
	case "json":
		// Inject into export struct
		export := ExportedData{
			Data: c.GetStore(),
		}

		// Marshal to JSON then print to stdout
		out, err := json.Marshal(export)
		if err != nil {
			return err
		}

		raw := make(map[string]any)
		if err := json.Unmarshal(out, &raw); err == nil {
			f := colorjson.NewFormatter()
			f.Indent = 2
			colorized, err := f.Marshal(raw)
			if err != nil {
				panic(err)
			}
			cmd.OutOrStdout().Write(colorized)
			return nil
		} else {
			cmd.OutOrStdout().Write(out)
			return nil
		}
	default:
		return r.Render(cmd, args)
	}
}
