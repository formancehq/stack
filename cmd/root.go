package cmd

import (
	"os"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedlogging/sharedlogginglogrus"
	"github.com/numary/webhooks-cloud/api/server"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "develop"

func Execute() {
	logger := logrus.New()
	loggerFactory := sharedlogging.StaticLoggerFactory(sharedlogginglogrus.New(logger))
	sharedlogging.SetFactory(loggerFactory)

	viper.SetDefault("version", Version)

	rootCmd := &cobra.Command{
		Use:  "webhooks",
		RunE: server.Start,
	}

	rootCmd.PersistentFlags().String(constants.ServerHttpBindAddressFlag,
		constants.DefaultBindAddress, "API bind address")
	rootCmd.PersistentFlags().String(constants.StorageMongoConnStringFlag,
		constants.DefaultMongoConnString, "Mongo connection string")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		sharedlogging.Errorf("viper.BindFlags: %s", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		sharedlogging.Errorf("cobra.Command.Execute: %s", err)
		os.Exit(1)
	}
}
