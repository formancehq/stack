package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/service"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use: "webhooks",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	otlptraces.InitOTLPTracesFlags(root.PersistentFlags())
	publish.InitCLIFlags(root)
	auth.InitAuthFlags(root.PersistentFlags())
	flag.Init(root.PersistentFlags())
	bunconnect.InitFlags(root.PersistentFlags())
	iam.InitFlags(root.PersistentFlags())
	service.BindFlags(root)
	licence.InitCLIFlags(root)

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
