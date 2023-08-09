package modules

import (
	"fmt"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
)

type monitoringType string

const (
	monitoringTypeTraces  monitoringType = "TRACES"
	monitoringTypeMetrics monitoringType = "METRICS"
)

func MonitoringTracesEnvVars(traces *stackv1beta3.TracesSpec, prefix string) ContainerEnv {
	if traces.Otlp != nil {
		return monitoringOTLPEnvVars(traces.Otlp, prefix, monitoringTypeTraces)
	}
	return ContainerEnv{}
}

func MonitoringMetricsEnvVars(metrics *stackv1beta3.MetricsSpec, prefix string) ContainerEnv {
	if metrics.Otlp != nil {
		return monitoringOTLPEnvVars(metrics.Otlp, prefix, monitoringTypeMetrics)
	}
	return ContainerEnv{}
}

func monitoringOTLPEnvVars(otlp *stackv1beta3.OtlpSpec, prefix string, monitoringType monitoringType) ContainerEnv {
	return ContainerEnv{
		Env(fmt.Sprintf("%sOTEL_%s", prefix, string(monitoringType)), "true"),
		Env(fmt.Sprintf("%sOTEL_%s_EXPORTER", prefix, string(monitoringType)), "otlp"),
		EnvFromBool(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_INSECURE", prefix, string(monitoringType)), otlp.Insecure),
		Env(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_MODE", prefix, string(monitoringType)), otlp.Mode),
		Env(fmt.Sprintf("%sOTEL_%s_PORT", prefix, string(monitoringType)), fmt.Sprint(otlp.Port)),
		Env(fmt.Sprintf("%sOTEL_%s_ENDPOINT", prefix, string(monitoringType)), otlp.Endpoint),
		Env(fmt.Sprintf("%sOTEL_%s_EXPORTER_OTLP_ENDPOINT", prefix, string(monitoringType)), controllerutils.ComputeEnvVar(prefix, "%s:%s", fmt.Sprintf("OTEL_%s_ENDPOINT", string(monitoringType)), fmt.Sprintf("OTEL_%s_PORT", string(monitoringType)))),
		Env(fmt.Sprintf("%sOTEL_RESOURCE_ATTRIBUTES", prefix), otlp.ResourceAttributes),
		Env(fmt.Sprintf("%sOTEL_EXPORTER_OTLP_TRACES_ENDPOINT", prefix), controllerutils.ComputeEnvVar(prefix, "http://%s", fmt.Sprintf("OTEL_TRACES_EXPORTER_OTLP_ENDPOINT"))),
	}
}
