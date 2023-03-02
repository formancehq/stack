package modules

import (
	"fmt"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
)

func MonitoringEnvVarsWithPrefix(m stackv1beta3.MonitoringSpec) ContainerEnv {
	if m.Traces != nil {
		return monitoringTracesEnvVars(m.Traces)
	}
	return ContainerEnv{}
}

func monitoringTracesEnvVars(traces *stackv1beta3.TracesSpec) ContainerEnv {
	if traces.Otlp != nil {
		return monitoringTracesOTLPEnvVars(traces.Otlp)
	}
	return ContainerEnv{}
}

func monitoringTracesOTLPEnvVars(otlp *stackv1beta3.TracesOtlpSpec) ContainerEnv {
	return ContainerEnv{
		Env("OTEL_TRACES", "true"),
		Env("OTEL_TRACES_EXPORTER", "otlp"),
		EnvFromBool("OTEL_TRACES_EXPORTER_OTLP_INSECURE", otlp.Insecure),
		Env("OTEL_TRACES_EXPORTER_OTLP_MODE", otlp.Mode),
		Env("OTEL_TRACES_PORT", fmt.Sprint(otlp.Port)),
		Env("OTEL_TRACES_ENDPOINT", otlp.Endpoint),
		Env("OTEL_TRACES_EXPORTER_OTLP_ENDPOINT", controllerutils.ComputeEnvVar("", "%s:%s", "OTEL_TRACES_ENDPOINT", "OTEL_TRACES_PORT")),
		Env("OTEL_RESOURCE_ATTRIBUTES", otlp.ResourceAttributes),
	}
}
