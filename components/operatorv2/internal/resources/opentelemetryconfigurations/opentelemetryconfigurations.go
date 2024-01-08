package opentelemetryconfigurations

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"k8s.io/api/core/v1"
)

type MonitoringType string

const (
	MonitoringTypeTraces  MonitoringType = "TRACES"
	MonitoringTypeMetrics MonitoringType = "METRICS"
)

func EnvVarsIfEnabled(ctx core.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := stacks.GetByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, stackName)
	if err != nil {
		return nil, err
	}
	if configuration != nil {
		return GetEnvVars(configuration, serviceName), nil
	}
	return nil, nil
}

func GetEnvVars(config *v1beta1.OpenTelemetryConfiguration, serviceName string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)
	if config.Spec.Traces != nil {
		if config.Spec.Traces.Otlp != nil {
			ret = append(ret, envVars(config.Spec.Traces.Otlp, MonitoringTypeTraces, serviceName)...)
		}
	}
	if config.Spec.Metrics != nil {
		if config.Spec.Metrics.Otlp != nil {
			ret = append(ret, envVars(config.Spec.Metrics.Otlp, MonitoringTypeMetrics, serviceName)...)
		}
	}
	return ret
}

func envVars(otlp *v1beta1.OtlpSpec, monitoringType MonitoringType, serviceName string) []v1.EnvVar {
	return []v1.EnvVar{
		core.Env(fmt.Sprintf("OTEL_%s", string(monitoringType)), "true"),
		core.Env(fmt.Sprintf("OTEL_%s_EXPORTER", string(monitoringType)), "otlp"),
		core.EnvFromBool(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_INSECURE", string(monitoringType)), otlp.Insecure),
		core.Env(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_MODE", string(monitoringType)), otlp.Mode),
		core.Env(fmt.Sprintf("OTEL_%s_PORT", string(monitoringType)), fmt.Sprint(otlp.Port)),
		core.Env(fmt.Sprintf("OTEL_%s_ENDPOINT", string(monitoringType)), otlp.Endpoint),
		core.Env(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_ENDPOINT", string(monitoringType)), core.ComputeEnvVar("%s:%s", fmt.Sprintf("OTEL_%s_ENDPOINT", string(monitoringType)), fmt.Sprintf("OTEL_%s_PORT", string(monitoringType)))),
		core.Env("OTEL_RESOURCE_ATTRIBUTES", otlp.ResourceAttributes),
		core.Env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", core.ComputeEnvVar("http://%s", fmt.Sprintf("OTEL_TRACES_EXPORTER_OTLP_ENDPOINT"))),
		core.Env("OTEL_SERVICE_NAME", serviceName),
	}
}