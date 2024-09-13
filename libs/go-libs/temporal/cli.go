package temporal

import (
	"github.com/spf13/pflag"
)

const (
	TemporalAddressFlag               = "temporal-address"
	TemporalNamespaceFlag             = "temporal-namespace"
	TemporalSSLClientKeyFlag          = "temporal-ssl-client-key"
	TemporalSSLClientCertFlag         = "temporal-ssl-client-cert"
	TemporalTaskQueueFlag             = "temporal-task-queue"
	TemporalInitSearchAttributesFlag  = "temporal-init-search-attributes"
	TemporalMaxParallelActivitiesFlag = "temporal-max-parallel-activities"
)

func AddFlags(flags *pflag.FlagSet) {
	flags.String(TemporalAddressFlag, "", "Temporal server address")
	flags.String(TemporalNamespaceFlag, "default", "Temporal namespace")
	flags.String(TemporalSSLClientKeyFlag, "", "Temporal client key")
	flags.String(TemporalSSLClientCertFlag, "", "Temporal client cert")
	flags.String(TemporalTaskQueueFlag, "default", "Temporal task queue name")
	flags.Bool(TemporalInitSearchAttributesFlag, false, "Init temporal search attributes")
	flags.Float64(TemporalMaxParallelActivitiesFlag, 10, "Maximum number of parallel activities")
}
