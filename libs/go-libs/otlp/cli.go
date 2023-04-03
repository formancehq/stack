package otlp

import (
	"sync"

	flag "github.com/spf13/pflag"
)

var (
	once sync.Once
)

const (
	OtelResourceAttributes             = "otel-resource-attributes"
	OtelServiceName                    = "otel-service-name"
	OtelTracesExporterOTLPModeFlag     = "otel-traces-exporter-otlp-mode"
	OtelTracesExporterOTLPEndpointFlag = "otel-traces-exporter-otlp-endpoint"
	OtelTracesExporterOTLPInsecureFlag = "otel-traces-exporter-otlp-insecure"
)

func InitOTLPTracesFlags(flags *flag.FlagSet) {
	once.Do(func() {
		flags.String(OtelServiceName, "", "OpenTelemetry service name")
		flags.StringSlice(OtelResourceAttributes, []string{}, "Additional OTLP resource attributes")
		flags.String(OtelTracesExporterOTLPModeFlag, "grpc", "OpenTelemetry traces OTLP exporter mode (grpc|http)")
		flags.String(OtelTracesExporterOTLPEndpointFlag, "", "OpenTelemetry traces grpc endpoint")
		flags.Bool(OtelTracesExporterOTLPInsecureFlag, false, "OpenTelemetry traces grpc insecure")
	})
}
