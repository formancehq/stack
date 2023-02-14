package modules

import (
	"fmt"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
)

func MonitoringEnvVarsWithPrefix(m stackv1beta3.MonitoringSpec, prefix string) ContainerEnv {
	if m.Traces != nil {
		return monitoringTracesEnvVars(m.Traces, prefix)
	}
	return ContainerEnv{}
}

func monitoringTracesEnvVars(traces *stackv1beta3.TracesSpec, prefix string) ContainerEnv {
	if traces.Otlp != nil {
		return monitoringTracesOTLPEnvVars(traces.Otlp, prefix)
	}
	return ContainerEnv{}
}

func monitoringTracesOTLPEnvVars(otlp *stackv1beta3.TracesOtlpSpec, prefix string) ContainerEnv {
	return ContainerEnv{
		Env(prefix+"OTEL_TRACES", "true"),
		Env(prefix+"OTEL_TRACES_EXPORTER", "otlp"),
		EnvFromBool(prefix+"OTEL_TRACES_EXPORTER_OTLP_INSECURE", otlp.Insecure),
		Env(prefix+"OTEL_TRACES_EXPORTER_OTLP_MODE", otlp.Mode),
		Env(prefix+"OTEL_TRACES_PORT", fmt.Sprint(otlp.Port)),
		Env(prefix+"OTEL_TRACES_ENDPOINT", otlp.Endpoint),
		Env(prefix+"OTEL_TRACES_EXPORTER_OTLP_ENDPOINT", controllerutils.ComputeEnvVar(prefix, "%s:%s", "OTEL_TRACES_ENDPOINT", "OTEL_TRACES_PORT")),
		Env(prefix+"OTEL_RESOURCE_ATTRIBUTES", otlp.ResourceAttributes),
	}
}
