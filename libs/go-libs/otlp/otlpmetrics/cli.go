package otlpmetrics

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/otlp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	OtelMetricsFlag                                   = "otel-metrics"
	OtelMetricsExporterPushIntervalFlag               = "otel-metrics-exporter-push-interval"
	OtelMetricsRuntimeFlag                            = "otel-metrics-runtime"
	OtelMetricsRuntimeMinimumReadMemStatsIntervalFlag = "otel-metrics-runtime-minimum-read-mem-stats-interval"
)

func InitOTLPTracesFlags(flags *flag.FlagSet) {
	otlp.InitOTLPTracesFlags(flags)

	flags.Bool(OtelMetricsFlag, false, "Enable OpenTelemetry traces support")
	flags.Duration(OtelMetricsExporterPushIntervalFlag, 100*time.Millisecond, "OpenTelemetry metrics exporter push interval")
	flags.Bool(OtelMetricsRuntimeFlag, false, "Enable OpenTelemetry runtime metrics")
	flags.Duration(OtelMetricsRuntimeMinimumReadMemStatsIntervalFlag, 15*time.Second, "OpenTelemetry runtime metrics minimum read mem stats interval")
}

func CLIMetricsModule(v *viper.Viper) fx.Option {
	if v.GetBool(OtelMetricsFlag) {
		return MetricsModule(ModuleConfig{
			ServiceName:    v.GetString(otlp.OtelServiceName),
			ServiceVersion: "develop",
			OTLPConfig: &OTLPConfig{
				Mode:     v.GetString(otlp.OtelTracesExporterOTLPModeFlag),
				Endpoint: v.GetString(otlp.OtelTracesExporterOTLPEndpointFlag),
				Insecure: v.GetBool(otlp.OtelTracesExporterOTLPInsecureFlag),
			},
			RuntimeMetrics:              v.GetBool(OtelMetricsRuntimeFlag),
			MinimumReadMemStatsInterval: v.GetDuration(OtelMetricsRuntimeMinimumReadMemStatsIntervalFlag),
			PushInterval:                v.GetDuration(OtelMetricsExporterPushIntervalFlag),
			ResourceAttributes:          v.GetStringSlice(otlp.OtelResourceAttributes),
		})
	}
	return fx.Options()
}
