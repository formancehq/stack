package common

import (
	"bytes"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/gateways"
	"github.com/formancehq/operator/v2/internal/opentelemetryconfigurations"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	utils2 "github.com/formancehq/operator/v2/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"text/template"
)

func GetDevEnvVars(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}) []v1.EnvVar {
	return []v1.EnvVar{
		utils2.EnvFromBool("DEBUG", stack.Spec.Debug || service.IsDebug()),
		utils2.EnvFromBool("DEV", stack.Spec.Dev || service.IsDev()),
		utils2.Env("STACK", stack.Name),
	}
}

func GetCommonServicesEnvVars(ctx reconcilers.Context, stack *v1beta1.Stack, serviceName string, service interface {
	IsDebug() bool
	IsDev() bool
}) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)
	configuration, err := opentelemetryconfigurations.GetOpenTelemetryConfiguration(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	if configuration != nil {
		ret = append(ret, opentelemetryconfigurations.MonitoringEnvVars(configuration, serviceName)...)
	}

	env, err := gateways.GetURLSAsEnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	ret = append(ret, env...)
	ret = append(ret, GetDevEnvVars(stack, service)...)

	return ret, nil
}

func StandardHTTPPort() v1.ContainerPort {
	return v1.ContainerPort{
		Name:          "http",
		ContainerPort: 8080,
	}
}

func GetVersion(stack *v1beta1.Stack, defaultVersion string) string {
	if defaultVersion == "" {
		return stack.GetVersion()
	}
	return defaultVersion
}

func GetStack(ctx reconcilers.Context, spec interface {
	GetStack() string
}) (*v1beta1.Stack, error) {
	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: spec.GetStack(),
	}, stack); err != nil {
		return nil, err
	}

	return stack, nil
}

func CreateCaddyfileConfigMap(ctx reconcilers.Context, stack *v1beta1.Stack,
	name, _tpl string, additionalData map[string]any, options ...utils2.ObjectMutator[*v1.ConfigMap]) (*v1.ConfigMap, error) {
	caddyfile, err := ComputeCaddyfile(ctx, stack, _tpl, additionalData)
	if err != nil {
		return nil, err
	}

	options = append([]utils2.ObjectMutator[*v1.ConfigMap]{
		func(t *v1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}
		},
	}, options...)

	configMap, _, err := utils2.CreateOrUpdate[*v1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      name,
	},
		options...,
	)
	return configMap, err
}

func ComputeCaddyfile(ctx reconcilers.Context, stack *v1beta1.Stack, _tpl string, additionalData map[string]any) (string, error) {
	tpl := template.Must(template.New("main").Parse(_tpl))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := opentelemetryconfigurations.IsOpenTelemetryEnabled(ctx, stack.Name)
	if err != nil {
		return "", err
	}

	data := map[string]any{
		"EnableOpenTelemetry": openTelemetryEnabled,
	}
	data = collectionutils.MergeMaps(data, additionalData)

	if err := tpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
