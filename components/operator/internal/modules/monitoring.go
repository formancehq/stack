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
		Env(fmt.Sprintf("%sOTEL_TRACES", prefix), "true"),
		Env(fmt.Sprintf("%sOTEL_TRACES_EXPORTER", prefix), "otlp"),
		EnvFromBool(fmt.Sprintf("%sOTEL_TRACES_EXPORTER_OTLP_INSECURE", prefix), otlp.Insecure),
		Env(fmt.Sprintf("%sOTEL_TRACES_EXPORTER_OTLP_MODE", prefix), otlp.Mode),
		Env(fmt.Sprintf("%sOTEL_TRACES_PORT", prefix), fmt.Sprint(otlp.Port)),
		Env(fmt.Sprintf("%sOTEL_TRACES_ENDPOINT", prefix), otlp.Endpoint),
		Env(fmt.Sprintf("%sOTEL_TRACES_EXPORTER_OTLP_ENDPOINT", prefix), controllerutils.ComputeEnvVar(prefix, "%s:%s", "OTEL_TRACES_ENDPOINT", "OTEL_TRACES_PORT")),
		Env(fmt.Sprintf("%sOTEL_RESOURCE_ATTRIBUTES", prefix), otlp.ResourceAttributes),
	}
}
