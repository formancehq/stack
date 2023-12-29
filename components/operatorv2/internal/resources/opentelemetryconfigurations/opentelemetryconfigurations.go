package opentelemetryconfigurations

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetOpenTelemetryConfiguration(ctx core.Context, stackName string) (*v1beta1.OpenTelemetryConfiguration, error) {
	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	openTelemetryTracesList := &v1beta1.OpenTelemetryConfigurationList{}
	if err := ctx.GetClient().List(ctx, openTelemetryTracesList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(openTelemetryTracesList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &openTelemetryTracesList.Items[0], nil
	default:
		return nil, errors.New("found multiple opentelemetry config")
	}
}

func IsOpenTelemetryEnabled(ctx core.Context, stackName string) (bool, error) {
	configuration, err := GetOpenTelemetryConfiguration(ctx, stackName)
	if err == nil {
		return false, err
	}
	if configuration == nil {
		return false, nil
	}
	return true, nil
}

type MonitoringType string

const (
	MonitoringTypeTraces  MonitoringType = "TRACES"
	MonitoringTypeMetrics MonitoringType = "METRICS"
)

func MonitoringEnvVars(config *v1beta1.OpenTelemetryConfiguration, serviceName string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)
	if config.Spec.Traces != nil {
		if config.Spec.Traces.Otlp != nil {
			ret = append(ret, MonitoringOTLPEnvVars(config.Spec.Traces.Otlp, MonitoringTypeTraces, serviceName)...)
		}
	}
	if config.Spec.Metrics != nil {
		if config.Spec.Metrics.Otlp != nil {
			ret = append(ret, MonitoringOTLPEnvVars(config.Spec.Metrics.Otlp, MonitoringTypeMetrics, serviceName)...)
		}
	}
	return nil
}

func MonitoringOTLPEnvVars(otlp *v1beta1.OtlpSpec, monitoringType MonitoringType, serviceName string) []v1.EnvVar {
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
