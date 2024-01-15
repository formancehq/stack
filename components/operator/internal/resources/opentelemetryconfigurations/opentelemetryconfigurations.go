package opentelemetryconfigurations

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

type MonitoringType string

const (
	MonitoringTypeTraces  MonitoringType = "TRACES"
	MonitoringTypeMetrics MonitoringType = "METRICS"
)

func EnvVarsIfEnabled(ctx core.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	return EnvVarsIfEnabledWithPrefix(ctx, stackName, serviceName, "")
}

func EnvVarsIfEnabledWithPrefix(ctx core.Context, stackName, serviceName, prefix string) ([]v1.EnvVar, error) {
	configuration, err := core.GetByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, stackName)
	if err != nil {
		return nil, err
	}
	if configuration != nil {
		return GetEnvVarsWithPrefix(configuration, serviceName, prefix), nil
	}
	return nil, nil
}

func GetEnvVars(config *v1beta1.OpenTelemetryConfiguration, serviceName string) []v1.EnvVar {
	return GetEnvVarsWithPrefix(config, serviceName, "")
}

func GetEnvVarsWithPrefix(config *v1beta1.OpenTelemetryConfiguration, serviceName, prefix string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)
	if config.Spec.Traces != nil {
		if config.Spec.Traces.Otlp != nil {
			ret = append(ret, envVars(config.Spec.Traces.Otlp, MonitoringTypeTraces, serviceName, prefix)...)
		}
	}
	if config.Spec.Metrics != nil {
		if config.Spec.Metrics.Otlp != nil {
			ret = append(ret, envVars(config.Spec.Metrics.Otlp, MonitoringTypeMetrics, serviceName, prefix)...)
		}
	}
	return ret
}

func envVars(otlp *v1beta1.OtlpSpec, monitoringType MonitoringType, serviceName, prefix string) []v1.EnvVar {

	ret := []v1.EnvVar{
		core.Env(fmt.Sprintf("%sOTEL_%s", prefix, string(monitoringType)), "true"),
		core.Env(fmt.Sprintf("%sOTEL_%s_EXPORTER", prefix, string(monitoringType)), "otlp"),
		core.EnvFromBool(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_INSECURE", prefix, string(monitoringType)), otlp.Insecure),
		core.Env(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_MODE", prefix, string(monitoringType)), otlp.Mode),
		core.Env(fmt.Sprintf("%sOTEL_%s_PORT", prefix, string(monitoringType)), fmt.Sprint(otlp.Port)),
		core.Env(fmt.Sprintf("%sOTEL_%s_ENDPOINT", prefix, string(monitoringType)), otlp.Endpoint),
		core.Env(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_ENDPOINT", prefix, string(monitoringType)), core.ComputeEnvVar("%s:%s",
			fmt.Sprintf("%sOTEL_%s_ENDPOINT", prefix, string(monitoringType)),
			fmt.Sprintf("%sOTEL_%s_PORT", prefix, string(monitoringType)))),
		core.Env(fmt.Sprintf("%sOTEL_EXPORTER_OTLP_TRACES_ENDPOINT", prefix),
			core.ComputeEnvVar("http://%s", fmt.Sprintf("%sOTEL_TRACES_EXPORTER_OTLP_ENDPOINT", prefix))),
		core.Env(fmt.Sprintf("%sOTEL_SERVICE_NAME", prefix), serviceName),
	}

	if otlp.ResourceAttributes != nil {
		resourceAttributes := ""
		for k, v := range otlp.ResourceAttributes {
			resourceAttributes = fmt.Sprintf("%s%s=%s ", resourceAttributes, k, v)
		}
		ret = append(ret, core.Env(fmt.Sprintf("%sOTEL_RESOURCE_ATTRIBUTES", prefix), resourceAttributes))
	}

	return ret
}
