package formance_com

import (
	"bytes"
	"fmt"
	"golang.org/x/mod/semver"
	"strings"
	"text/template"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
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
	return GetDevEnvVarsWithPrefix(stack, service, "")
}

func GetDevEnvVarsWithPrefix(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}, prefix string) []v1.EnvVar {
	return []v1.EnvVar{
		EnvFromBool(fmt.Sprintf("%sDEBUG", prefix), stack.Spec.Debug || service.IsDebug()),
		EnvFromBool(fmt.Sprintf("%sDEV", prefix), stack.Spec.Dev || service.IsDev()),
		Env(fmt.Sprintf("%sSTACK", prefix), stack.Name),
	}
}

func ComputeCaddyfile(ctx Context, stack *v1beta1.Stack, _tpl string, additionalData map[string]any) (string, error) {
	tpl := template.Must(template.New("main").Funcs(map[string]any{
		"join":            strings.Join,
		"semver_compare":  semver.Compare,
		"semver_is_valid": semver.IsValid,
	}).Parse(_tpl))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := IsEnabledByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, stack.Name)
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
