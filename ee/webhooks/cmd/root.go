package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use: "webhooks",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding viper flags")
		},
	}

	otlptraces.InitOTLPTracesFlags(root.PersistentFlags())
	publish.InitCLIFlags(root)
	auth.InitAuthFlags(root.PersistentFlags())
	flag.Init(root.PersistentFlags())

	root.AddCommand(newServeCommand())
	root.AddCommand(newWorkerCommand())
	root.AddCommand(newVersionCommand())
	root.AddCommand(newMigrateCommand())

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		logging.Errorf("cobra.Command.Execute: %s", err)
		os.Exit(1)
	}
}
