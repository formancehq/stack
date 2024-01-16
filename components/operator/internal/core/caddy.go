package core

import (
	"bytes"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"text/template"
)

func ConfigureCaddy(caddyfile *v1.ConfigMap, image string, env []v1.EnvVar,
	resourceRequirements *v1.ResourceRequirements) ObjectMutator[*appsv1.Deployment] {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Annotations = collectionutils.MergeMaps(t.Spec.Template.Annotations, map[string]string{
			"caddyfile-hash": HashFromConfigMaps(caddyfile),
		})
		t.Spec.Template.Spec.Volumes = []v1.Volume{
			volumeFromConfigMap("caddyfile", caddyfile),
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
				Image:     image,
				Env:       env,
				Resources: GetResourcesRequirementsWithDefault(resourceRequirements, ResourceSizeSmall()),
				VolumeMounts: []v1.VolumeMount{
					NewVolumeMount("caddyfile", "/gateway"),
				},
				Ports: []v1.ContainerPort{{
					Name:          "http",
					ContainerPort: 8080,
				}},
			},
		}
	}
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
