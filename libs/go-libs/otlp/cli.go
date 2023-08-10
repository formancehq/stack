package otlp

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

const (
	OtelResourceAttributes = "otel-resource-attributes"
	OtelServiceName        = "otel-service-name"
)

func InitOTLPFlags(flags *flag.FlagSet) {
	flags.String(OtelServiceName, "", "OpenTelemetry service name")
	flags.StringSlice(OtelResourceAttributes, []string{}, "Additional OTLP resource attributes")
}

func LoadResource(serviceName string, resourceAttributes []string) fx.Option {
	return fx.Provide(func() (*resource.Resource, error) {
		defaultResource := resource.Default()
		attributes := make([]attribute.KeyValue, 0)
		if serviceName != "" {
			attributes = append(attributes, attribute.String("service.name", serviceName))
		}
		for _, ra := range resourceAttributes {
			parts := strings.SplitN(ra, "=", 2)
			if len(parts) < 2 {
				return nil, fmt.Errorf("malformed otlp attribute: %s", ra)
			}
			attributes = append(attributes, attribute.String(parts[0], parts[1]))
		}
		return resource.Merge(defaultResource, resource.NewSchemaless(attributes...))
	})
}
