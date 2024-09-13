package caddy

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/internal/core"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func DeploymentTemplate(ctx core.Context, stack *v1beta1.Stack, owner v1beta1.Object, caddyfile *v1.ConfigMap, image string, env []v1.EnvVar) (*appsv1.Deployment, error) {
	t := &appsv1.Deployment{}

	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, owner))
	if err != nil {
		return nil, err
	}

	if len(otlpEnv) > 0 {
		env = append(env, otlpEnv...)
		scheme := "https"
		for _, envVar := range env {
			if envVar.Name == "OTEL_TRACES_EXPORTER_OTLP_INSECURE" {
				if envVar.Value == "true" {
					scheme = "http"
				}
				break
			}
		}
		env = append(env, core.Env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", fmt.Sprintf("%s://$(OTEL_TRACES_EXPORTER_OTLP_ENDPOINT)", scheme)))
		env = append(env, core.Env("OTEL_EXPORTER_OTLP_PROTOCOL", "$(OTEL_TRACES_EXPORTER_OTLP_MODE)"))
	}

	t.Spec.Template.Annotations = collectionutils.MergeMaps(t.Spec.Template.Annotations, map[string]string{
		"caddyfile-hash": core.HashFromConfigMaps(caddyfile),
	})
	t.Spec.Template.Spec.Volumes = []v1.Volume{
		volumeFromConfigMap("caddyfile", caddyfile),
		{
			Name: "data",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}
	t.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:    "gateway",
			Command: []string{"/usr/bin/caddy"},
			Args: []string{
				"run",
				"--config", "/gateway/Caddyfile",
				"--adapter", "caddyfile",
			},
			Image: image,
			Env:   env,
			VolumeMounts: []v1.VolumeMount{
				core.NewVolumeMount("caddyfile", "/gateway", true),
				core.NewVolumeMount("data", "/data", false),
			},
			Ports: []v1.ContainerPort{{
				Name:          "http",
				ContainerPort: 8080,
			}, {
				Name:          "metrics",
				ContainerPort: 3080,
			}},
			SecurityContext: &v1.SecurityContext{
				Capabilities: &v1.Capabilities{
					Add: []v1.Capability{"NET_BIND_SERVICE"},
				},
			},
		},
	}

	return t, nil
}

func volumeFromConfigMap(name string, cm *v1.ConfigMap) v1.Volume {
	return v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: cm.Name,
				},
			},
		},
	}
}

func ComputeCaddyfile(ctx core.Context, stack *v1beta1.Stack, _tpl string, additionalData map[string]any) (string, error) {
	tpl := template.Must(template.New("main").Funcs(map[string]any{
		"join":            strings.Join,
		"semver_compare":  semver.Compare,
		"semver_is_valid": semver.IsValid,
	}).Parse(_tpl))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := settings.HasOpenTelemetryTracesEnabled(ctx, stack.Name)
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

func CreateCaddyfileConfigMap(ctx core.Context, stack *v1beta1.Stack,
	name, _tpl string, additionalData map[string]any, options ...core.ObjectMutator[*v1.ConfigMap]) (*v1.ConfigMap, error) {
	caddyfile, err := ComputeCaddyfile(ctx, stack, _tpl, additionalData)
	if err != nil {
		return nil, err
	}

	options = append([]core.ObjectMutator[*v1.ConfigMap]{
		func(t *v1.ConfigMap) error {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}

			return nil
		},
	}, options...)

	configMap, _, err := core.CreateOrUpdate[*v1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      name,
	},
		options...,
	)
	return configMap, err
}
