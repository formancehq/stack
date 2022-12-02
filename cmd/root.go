package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/spf13/cobra"
)

var retriesSchedule []time.Duration

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use: "webhooks",
	}

	var err error
	sharedotlptraces.InitOTLPTracesFlags(root.PersistentFlags())
	retriesSchedule, err = flag.Init(root.PersistentFlags())
	cobra.CheckErr(err)

	root.AddCommand(serverCmd)
	root.AddCommand(workerCmd)
	root.AddCommand(versionCmd)

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		sharedlogging.Errorf("cobra.Command.Execute: %s", err)
		os.Exit(1)
	}
}
