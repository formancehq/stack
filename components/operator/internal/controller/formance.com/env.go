package formance_com

import (
	"bytes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"text/template"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/stacks"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// TODO: The stack reconciler can create a config map container env var for dev and debug
// This way, we avoid the need to fetch the stack object at each reconciliation loop
func GetDevEnvVars(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}) []v1.EnvVar {
	return []v1.EnvVar{
		EnvFromBool("DEBUG", stack.Spec.Debug || service.IsDebug()),
		EnvFromBool("DEV", stack.Spec.Dev || service.IsDev()),
		Env("STACK", stack.Name),
	}
}

func GetCommonServicesEnvVars(ctx Context, stack *v1beta1.Stack, service interface {
	client.Object
	IsDebug() bool
	IsDev() bool
}) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)
	env, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, strings.ToLower(service.GetObjectKind().GroupVersionKind().Kind))
	if err != nil {
		return nil, err
	}
	ret = append(ret, env...)

	env, err = gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	ret = append(ret, env...)
	ret = append(ret, GetDevEnvVars(stack, service)...)

	return ret, nil
}

func ComputeCaddyfile(ctx Context, stack *v1beta1.Stack, _tpl string, additionalData map[string]any) (string, error) {
	tpl := template.Must(template.New("main").Funcs(map[string]any{
		"join": strings.Join,
	}).Parse(_tpl))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := stacks.IsEnabledByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, stack.Name)
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

func CreateCaddyfileConfigMap(ctx Context, stack *v1beta1.Stack,
	name, _tpl string, additionalData map[string]any, options ...ObjectMutator[*v1.ConfigMap]) (*v1.ConfigMap, error) {
	caddyfile, err := ComputeCaddyfile(ctx, stack, _tpl, additionalData)
	if err != nil {
		return nil, err
	}

	options = append([]ObjectMutator[*v1.ConfigMap]{
		func(t *v1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}
		},
	}, options...)

	configMap, _, err := CreateOrUpdate[*v1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      name,
	},
		options...,
	)
	return configMap, err
}
