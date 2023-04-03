package otlpmetrics

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/otlp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	OtelMetricsFlag                     = "otel-metrics"
	OtelMetricsExporterPushIntervalFlag = "otel-metrics-exporter-push-interval"
)

func InitOTLPTracesFlags(flags *flag.FlagSet) {
	otlp.InitOTLPTracesFlags(flags)

	flags.Bool(OtelMetricsFlag, false, "Enable OpenTelemetry traces support")
	flags.Duration(OtelMetricsExporterPushIntervalFlag, 100*time.Millisecond, "OpenTelemetry metrics exporter push interval")
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
			PushInterval:       v.GetDuration(OtelMetricsExporterPushIntervalFlag),
			ResourceAttributes: v.GetStringSlice(otlp.OtelResourceAttributes),
		})
	}
	return fx.Options()
}
