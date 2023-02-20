package cmd

import (
	"fmt"
	"os"

	_ "github.com/formancehq/orchestration/internal/workflow/stages/all"
	"github.com/formancehq/stack/libs/go-libs/app"
	"github.com/spf13/cobra"
)

var (
	ServiceName = "orchestration"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	stackURLFlag              = "stack-url"
	stackClientIDFlag         = "stack-client-id"
	stackClientSecretFlag     = "stack-client-secret"
	temporalAddressFlag       = "temporal-address"
	temporalNamespaceFlag     = "temporal-namespace"
	temporalSSLClientKeyFlag  = "temporal-ssl-client-key"
	temporalSSLClientCertFlag = "temporal-ssl-client-cert"
	temporalTaskQueueFlag     = "temporal-task-queue"
	postgresDSNFlag           = "postgres-dsn"
)

var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithCode(1, err)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP(app.DebugFlag, "d", false, "Debug mode")
	rootCmd.PersistentFlags().String(stackURLFlag, "", "Stack url")
	rootCmd.PersistentFlags().String(stackClientIDFlag, "", "Stack client ID")
	rootCmd.PersistentFlags().String(stackClientSecretFlag, "", "Stack client secret")
	rootCmd.PersistentFlags().String(temporalAddressFlag, "", "Temporal server address")
	rootCmd.PersistentFlags().String(temporalNamespaceFlag, "default", "Temporal namespace")
	rootCmd.PersistentFlags().String(temporalSSLClientKeyFlag, "", "Temporal client key")
	rootCmd.PersistentFlags().String(temporalSSLClientCertFlag, "", "Temporal client cert")
	rootCmd.PersistentFlags().String(temporalTaskQueueFlag, "default", "Temporal task queue name")
	rootCmd.PersistentFlags().String(postgresDSNFlag, "", "Postgres address")
}
