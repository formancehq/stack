package temporal

import (
	"github.com/spf13/cobra"
)

const (
	TemporalAddressFlag           = "temporal-address"
	TemporalNamespaceFlag         = "temporal-namespace"
	TemporalSSLClientKeyFlag      = "temporal-ssl-client-key"
	TemporalSSLClientCertFlag     = "temporal-ssl-client-cert"
	TemporalTaskQueueFlag         = "temporal-task-queue"
	TemporalInitSearchAttributes  = "temporal-init-search-attributes"
	TemporalMaxParallelActivities = "temporal-max-parallel-activities"
)

func InitCLIFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(TemporalAddressFlag, "", "Temporal server address")
	cmd.PersistentFlags().String(TemporalNamespaceFlag, "default", "Temporal namespace")
	cmd.PersistentFlags().String(TemporalSSLClientKeyFlag, "", "Temporal client key")
	cmd.PersistentFlags().String(TemporalSSLClientCertFlag, "", "Temporal client cert")
	cmd.PersistentFlags().String(TemporalTaskQueueFlag, "default", "Temporal task queue name")
	cmd.PersistentFlags().Bool(TemporalInitSearchAttributes, false, "Init temporal search attributes")
	cmd.PersistentFlags().Float64(TemporalMaxParallelActivities, 10, "Maximum number of parallel activities")
}
