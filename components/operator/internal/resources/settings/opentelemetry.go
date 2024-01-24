package settings

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

type MonitoringType string

const (
	MonitoringTypeTraces  MonitoringType = "TRACES"
	MonitoringTypeMetrics MonitoringType = "METRICS"
)

func GetOTELEnvVarsIfEnabled(ctx core.Context, stack *v1beta1.Stack, serviceName string) ([]v1.EnvVar, error) {
	return GetOTELEnvVarsIfEnabledWithPrefix(ctx, stack, serviceName, "")
}

func GetOTELEnvVarsIfEnabledWithPrefix(ctx core.Context, stack *v1beta1.Stack, serviceName, prefix string) ([]v1.EnvVar, error) {
	configuration, err := FindOpenTelemetryConfiguration(ctx, stack)
	if err != nil {
		return nil, err
	}
	if configuration != nil {
		return GetOTELEnvVarsWithPrefix(configuration, serviceName, prefix), nil
	}
	return nil, nil
}

func GetOTELEnvVars(config *v1beta1.OpenTelemetryConfiguration, serviceName string) []v1.EnvVar {
	return GetOTELEnvVarsWithPrefix(config, serviceName, "")
}

func GetOTELEnvVarsWithPrefix(config *v1beta1.OpenTelemetryConfiguration, serviceName, prefix string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)
	if config.Traces != nil {
		if config.Traces.Otlp != nil {
			ret = append(ret, otelEnvVars(config.Traces.Otlp, MonitoringTypeTraces, serviceName, prefix)...)
		}
	}
	if config.Metrics != nil {
		if config.Metrics.Otlp != nil {
			ret = append(ret, otelEnvVars(config.Metrics.Otlp, MonitoringTypeMetrics, serviceName, prefix)...)
		}
	}
	return ret
}

func otelEnvVars(otlp *v1beta1.OtlpSpec, monitoringType MonitoringType, serviceName, prefix string) []v1.EnvVar {

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

func FindOpenTelemetryConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.OpenTelemetryConfiguration, error) {
	tracesSpec, err := resolveTracesSpec(ctx, stack)
	if err != nil {
		return nil, err
	}
	metricsSpec, err := resolveMetricsSpec(ctx, stack)
	if err != nil {
		return nil, err
	}
	if tracesSpec == nil && metricsSpec == nil {
		return nil, nil
	}
	return &v1beta1.OpenTelemetryConfiguration{
		Traces:  tracesSpec,
		Metrics: metricsSpec,
	}, nil
}

func resolveMetricsSpec(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.MetricsSpec, error) {
	otlpSpec, err := resolveOTLPSpec(ctx, stack, "metrics")
	if err != nil {
		return nil, err
	}
	if otlpSpec == nil {
		return nil, nil
	}
	return &v1beta1.MetricsSpec{
		Otlp: otlpSpec,
	}, nil
}

func resolveTracesSpec(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.TracesSpec, error) {
	otlpSpec, err := resolveOTLPSpec(ctx, stack, "traces")
	if err != nil {
		return nil, err
	}
	if otlpSpec == nil {
		return nil, nil
	}
	return &v1beta1.TracesSpec{
		Otlp: otlpSpec,
	}, nil
}

func resolveOTLPSpec(ctx core.Context, stack *v1beta1.Stack, discr string) (*v1beta1.OtlpSpec, error) {

	dsn, err := GetStringOrEmpty(ctx, stack.Name, "opentelemetry", discr, "dsn")
	if err != nil {
		return nil, err
	}
	if dsn == "" {
		return nil, err
	}

	otlpURI, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if otlpURI.Scheme != "http" && otlpURI.Scheme != "grpc" {
		return nil, fmt.Errorf("invalid uri: %s", dsn)
	}

	port := uint64(4317)
	if otlpURI.Port() != "" {
		port, err = strconv.ParseUint(otlpURI.Port(), 10, 16)
		if err != nil {
			return nil, err
		}
	}

	resourceAttributes, err := GetMapOrEmpty(ctx, stack.Name, "opentelemetry", discr, "resource-attributes")
	if err != nil {
		return nil, err
	}

	return &v1beta1.OtlpSpec{
		Endpoint:           otlpURI.Hostname(),
		Port:               int32(port),
		Insecure:           isTrue(otlpURI.Query().Get("insecure")),
		Mode:               otlpURI.Scheme,
		ResourceAttributes: resourceAttributes,
	}, nil
}
