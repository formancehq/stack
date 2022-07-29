package cmd

import (
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var fxOptions = []fx.Option{
	// TODO: Add fx options
	fx.Invoke(func() {
		sharedlogging.Infof("App started.")
	}),
	fx.NopLogger,
}

var serveCmd = &cobra.Command{
	Use: "serve",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(fxOptions...)
		err := app.Start(cmd.Context())
		if err != nil {
			return err
		}
		<-app.Done()
		return app.Err()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
