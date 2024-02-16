/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package benthos

import (
	"embed"
	"fmt"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"sort"

	"github.com/formancehq/operator/internal/resources/services"
	"github.com/formancehq/operator/internal/resources/settings"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/deployments"
	benthosOperator "github.com/formancehq/operator/internal/resources/searches/benthos"
	"github.com/formancehq/search/benthos"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//+kubebuilder:rbac:groups=formance.com,resources=benthos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=benthos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=benthos/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, b *v1beta1.Benthos) error {

	if err := createDeployment(ctx, stack, b); err != nil {
		return err
	}

	if err := createService(ctx, b); err != nil {
		return err
	}

	return nil
}

func createService(ctx Context, b *v1beta1.Benthos) error {
	_, err := services.Create(ctx, b, "benthos", func(t *corev1.Service) error {
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": "benthos",
		}
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       4195,
				Protocol:   "TCP",
				TargetPort: intstr.FromString("http"),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": "benthos",
			},
		}

		return nil
	})
	return err
}

// TODO(gfyrag): there is a ton of search related configuration
// We need to this controller and keep it focused on benthos
func createDeployment(ctx Context, stack *v1beta1.Stack, b *v1beta1.Benthos) error {
	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return errors.Wrap(err, "searching broker configuration")
	}

	elasticSearchURI, err := settings.RequireURL(ctx, stack.Name, "elasticsearch.dsn")
	if err != nil {
		return errors.Wrap(err, "searching elasticsearch configuration")
	}

	var resourceReference *v1beta1.ResourceReference
	if secret := elasticSearchURI.Query().Get("secret"); secret != "" {
		resourceReference, err = resourcereferences.Create(ctx, b, "elasticsearch", secret, &corev1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, b, "elasticsearch")
	}
	if err != nil {
		return err
	}

	env := []corev1.EnvVar{
		Env("OPENSEARCH_URL", elasticSearchURI.WithoutQuery().String()),
		Env("TOPIC_PREFIX", b.Spec.Stack+"-"),
		Env("OPENSEARCH_INDEX", "stacks"),
		Env("STACK", b.Spec.Stack),
	}
	if b.Spec.Batching != nil {
		if b.Spec.Batching.Count != 0 {
			env = append(env, Env("OPENSEARCH_BATCHING_COUNT", fmt.Sprint(b.Spec.Batching.Count)))
		}
		if b.Spec.Batching.Period != "" {
			env = append(env, Env("OPENSEARCH_BATCHING_PERIOD", b.Spec.Batching.Period))
		}
	}

	otelEnvVars, err := settings.GetOTELEnvVars(ctx, stack.Name, "benthos")
	if err != nil {
		return err
	}
	env = append(env, otelEnvVars...)

	if brokerURI.Scheme == "kafka" {
		env = append(env, Env("KAFKA_ADDRESS", brokerURI.Host))
		if settings.IsTrue(brokerURI.Query().Get("tls")) {
			env = append(env, Env("KAFKA_TLS_ENABLED", "true"))
		}
		if settings.IsTrue(brokerURI.Query().Get("saslEnabled")) {
			env = append(env,
				Env("KAFKA_SASL_USERNAME", brokerURI.Query().Get("saslUsername")),
				Env("KAFKA_SASL_PASSWORD", brokerURI.Query().Get("saslPassword")),
				Env("KAFKA_SASL_MECHANISM", brokerURI.Query().Get("saslMechanism")),
			)
		}
	}
	if brokerURI.Scheme == "nats" {
		env = append(env, Env("NATS_URL", brokerURI.Host))
	}
	if secret := elasticSearchURI.Query().Get("secret"); elasticSearchURI.User != nil || secret != "" {
		env = append(env, Env("BASIC_AUTH_ENABLED", "true"))
		if secret == "" {
			password, _ := brokerURI.User.Password()
			env = append(env,
				Env("BASIC_AUTH_USERNAME", brokerURI.User.Username()),
				Env("BASIC_AUTH_PASSWORD", password),
			)
		} else {
			env = append(env,
				EnvFromSecret("BASIC_AUTH_USERNAME", secret, "username"),
				EnvFromSecret("BASIC_AUTH_PASSWORD", secret, "password"),
			)
		}
	} else {
		// Even if basic auth is not enabled, we need to set the env vars
		// to avoid benthos to crash due to linting errors
		env = append(env,
			Env("BASIC_AUTH_ENABLED", "false"),
			Env("BASIC_AUTH_USERNAME", "username"),
			Env("BASIC_AUTH_PASSWORD", "password"),
		)
	}

	cmd := []string{
		"/benthos",
		"-r", "/resources/*.yaml",
		"-t", "/templates/*.yaml",
	}

	openTelemetryEnabled, err := settings.HasOpenTelemetryTracesEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if openTelemetryEnabled {
		cmd = append(cmd, "-c", "/global/config.yaml")
	}

	cmd = append(cmd, "--log.level", "trace", "streams", "/streams/*.yaml")

	volumes := make([]corev1.Volume, 0)
	volumeMounts := make([]corev1.VolumeMount, 0)

	type directory struct {
		name string
		fs   embed.FS
	}

	directories := []directory{
		{
			name: "templates",
			fs:   benthos.Templates,
		},
		{
			name: "resources",
			fs:   benthos.Resources,
		},
	}
	if openTelemetryEnabled {
		directories = append(directories, directory{
			name: "global",
			fs:   benthosOperator.Global,
		})
	}

	if stack.Spec.EnableAudit {
		directories = append(directories, directory{
			name: "audit",
			fs:   benthosOperator.Audit,
		})
	} else {
		kinds, _, err := ctx.GetScheme().ObjectKinds(&corev1.ConfigMap{})
		if err != nil {
			return err
		}

		object := &unstructured.Unstructured{}
		object.SetGroupVersionKind(kinds[0])
		object.SetNamespace(stack.Name)
		object.SetName("benthos-audit")
		if err := client.IgnoreNotFound(ctx.GetClient().Delete(ctx, object)); err != nil {
			return errors.Wrap(err, "deleting audit config map")
		}
	}

	configMaps := make([]*corev1.ConfigMap, 0)

	for _, x := range directories {
		data := make(map[string]string)

		CopyDir(x.fs, x.name, x.name, &data)

		configMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
			Namespace: b.Spec.Stack,
			Name:      "benthos-" + x.name,
		},
			func(t *corev1.ConfigMap) error {
				t.Data = data

				return nil
			},
			WithController[*corev1.ConfigMap](ctx.GetScheme(), b),
		)
		if err != nil {
			return err
		}

		configMaps = append(configMaps, configMap)

		volumes = append(volumes, corev1.Volume{
			Name: x.name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "benthos-" + x.name,
					},
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      x.name,
			ReadOnly:  true,
			MountPath: "/" + x.name,
		})
	}

	if stack.Spec.EnableAudit {
		cmd = append(cmd, "/audit/gateway_audit.yaml")
	}

	streamList := &v1beta1.BenthosStreamList{}
	if err := ctx.GetClient().List(ctx, streamList, client.MatchingFields{
		"stack": b.Spec.Stack,
	}); err != nil {
		return err
	}

	streams := streamList.Items
	sort.Slice(streams, func(i, j int) bool {
		return streams[i].Name < streams[j].Name
	})

	_, err = deployments.CreateOrUpdate(ctx, stack, b, "benthos",
		resourcereferences.Annotate[*appsv1.Deployment]("elasticsearch-secret-hash", resourceReference),
		deployments.WithMatchingLabels("benthos"),
		deployments.WithInitContainers(b.Spec.InitContainers...),
		deployments.WithContainers(corev1.Container{
			Name:    "benthos",
			Image:   "public.ecr.aws/formance-internal/jeffail/benthos:v4.23.1-es",
			Env:     env,
			Command: cmd,
			Ports: []corev1.ContainerPort{{
				Name:          "http",
				ContainerPort: 4195,
			}},
			VolumeMounts: append(volumeMounts, corev1.VolumeMount{
				Name:      "streams",
				ReadOnly:  true,
				MountPath: "/streams",
			}),
		}),
		deployments.WithVolumes(append(volumes, corev1.Volume{
			Name: "streams",
			VolumeSource: corev1.VolumeSource{
				Projected: &corev1.ProjectedVolumeSource{
					Sources: Map(streams, func(stream v1beta1.BenthosStream) corev1.VolumeProjection {
						return corev1.VolumeProjection{
							ConfigMap: &corev1.ConfigMapProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: fmt.Sprintf("stream-%s", stream.Name),
								},
								Items: []corev1.KeyToPath{{
									Key:  "stream.yaml",
									Path: stream.Spec.Name + ".yaml",
								}},
							},
						}
					}),
				},
			},
		})...),
		func(t *appsv1.Deployment) error {
			t.Spec.Template.Annotations = MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"config-hash": HashFromConfigMaps(configMaps...),
			})

			return nil
		},
	)

	return err
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Benthos](),
			WithWatchDependency[*v1beta1.Benthos](&v1beta1.BenthosStream{}),
			WithOwn[*v1beta1.Benthos](&corev1.ConfigMap{}),
			WithOwn[*v1beta1.Benthos](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Benthos](&v1beta1.ResourceReference{}),
		),
	)
}
