package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "webhooks",
	}
	retriesSchedule []time.Duration
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		sharedlogging.Errorf("cobra.Command.Execute: %s", err)
		os.Exit(1)
	}
}

func init() {
	var err error
	retriesSchedule, err = flag.Init(rootCmd.PersistentFlags())
	cobra.CheckErr(err)
}
