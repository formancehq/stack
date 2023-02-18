package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/logging/logginglogrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

const (
	debugFlag = "debug"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := bindFlagsToViper(cmd); err != nil {
				return err
			}

			logrusLogger := logrus.New()
			if viper.GetBool(debugFlag) {
				logrusLogger.SetLevel(logrus.DebugLevel)
				logrusLogger.Infof("Debug mode enabled.")
			}
			logger := logginglogrus.New(logrusLogger)
			logging.SetFactory(logging.StaticLoggerFactory(logger))

			return nil
		},
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.PersistentFlags().BoolP(debugFlag, "d", false, "Debug mode")
	cmd.AddCommand(newServeCommand(), newVersionCommand())

	return cmd
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		exitWithCode(1, err)
	}
}
