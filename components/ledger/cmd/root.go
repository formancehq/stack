package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/ledger/cmd/internal"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	bindFlag = "bind"
)

var (
	Version                = "develop"
	BuildDate              = "-"
	Commit                 = "-"
	DefaultSegmentWriteKey = ""
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "ledger",
		Short:             "ledger",
		DisableAutoGenTag: true,
	}

	serve := NewServe()
	version := NewVersion()

	conf := NewConfig()
	conf.AddCommand(NewConfigInit())
	store := NewBucket()
	store.AddCommand(NewBucketInit())
	store.AddCommand(NewBucketList())
	store.AddCommand(NewBucketUpgrade())
	store.AddCommand(NewBucketUpgradeAll())
	store.AddCommand(NewBucketDelete())

	root.AddCommand(serve)
	root.AddCommand(conf)
	root.AddCommand(store)
	root.AddCommand(version)

	root.AddCommand(NewDocCommand())

	root.PersistentFlags().Bool(service.DebugFlag, false, "Debug mode")
	root.PersistentFlags().Bool(service.JsonFormattingLoggerFlag, true, "Json formatting mode for logger")
	root.PersistentFlags().String(bindFlag, "0.0.0.0:3068", "API bind address")

	otlpmetrics.InitOTLPMetricsFlags(root.PersistentFlags())
	otlptraces.InitOTLPTracesFlags(root.PersistentFlags())
	internal.InitAnalyticsFlags(root, DefaultSegmentWriteKey)
	publish.InitCLIFlags(root)
	driver.InitCLIFlags(root)

	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		panic(err)
	}

	internal.BindEnv(viper.GetViper())

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
