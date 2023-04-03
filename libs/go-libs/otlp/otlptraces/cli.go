package otlptraces

import (
	"github.com/formancehq/stack/libs/go-libs/otlp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	OtelTracesFlag                       = "otel-traces"
	OtelTracesBatchFlag                  = "otel-traces-batch"
	OtelTracesExporterFlag               = "otel-traces-exporter"
	OtelTracesExporterJaegerEndpointFlag = "otel-traces-exporter-jaeger-endpoint"
	OtelTracesExporterJaegerUserFlag     = "otel-traces-exporter-jaeger-user"
	OtelTracesExporterJaegerPasswordFlag = "otel-traces-exporter-jaeger-password"
)

func InitOTLPTracesFlags(flags *flag.FlagSet) {
	otlp.InitOTLPTracesFlags(flags)

	flags.Bool(OtelTracesFlag, false, "Enable OpenTelemetry traces support")
	flags.Bool(OtelTracesBatchFlag, false, "Use OpenTelemetry batching")
	flags.String(OtelTracesExporterFlag, "stdout", "OpenTelemetry traces exporter")
	flags.String(OtelTracesExporterJaegerEndpointFlag, "", "OpenTelemetry traces Jaeger exporter endpoint")
	flags.String(OtelTracesExporterJaegerUserFlag, "", "OpenTelemetry traces Jaeger exporter user")
	flags.String(OtelTracesExporterJaegerPasswordFlag, "", "OpenTelemetry traces Jaeger exporter password")
}

func CLITracesModule(v *viper.Viper) fx.Option {
	if v.GetBool(OtelTracesFlag) {
		return TracesModule(ModuleConfig{
			Batch:    v.GetBool(OtelTracesBatchFlag),
			Exporter: v.GetString(OtelTracesExporterFlag),
			JaegerConfig: func() *JaegerConfig {
				if v.GetString(OtelTracesExporterFlag) != JaegerExporter {
					return nil
				}
				return &JaegerConfig{
					Endpoint: v.GetString(OtelTracesExporterJaegerEndpointFlag),
					User:     v.GetString(OtelTracesExporterJaegerUserFlag),
					Password: v.GetString(OtelTracesExporterJaegerPasswordFlag),
				}
			}(),
			OTLPConfig: func() *OTLPConfig {
				if v.GetString(OtelTracesExporterFlag) != OTLPExporter {
					return nil
				}
				return &OTLPConfig{
					Mode:     v.GetString(otlp.OtelTracesExporterOTLPModeFlag),
					Endpoint: v.GetString(otlp.OtelTracesExporterOTLPEndpointFlag),
					Insecure: v.GetBool(otlp.OtelTracesExporterOTLPInsecureFlag),
				}
			}(),
			ServiceName:        v.GetString(otlp.OtelServiceName),
			ResourceAttributes: v.GetStringSlice(otlp.OtelResourceAttributes),
		})
	}
	return fx.Options()
}
