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

package streamprocessors

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
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

//+kubebuilder:rbac:groups=formance.com,resources=streamprocessors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=streamprocessors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=streamprocessors/finalizers,verbs=update

func Reconcile(ctx Context, streamProcessor *v1beta1.StreamProcessor) error {

	brokerConfiguration, err := RequireLabelledConfig[*v1beta1.BrokerConfiguration](ctx, streamProcessor.Spec.Stack)
	if err != nil {
		return errors.Wrap(err, "searching broker configuration")
	}

	elasticSearchConfiguration, err := RequireLabelledConfig[*v1beta1.ElasticSearchConfiguration](ctx, streamProcessor.Spec.Stack)
	if err != nil {
		return errors.Wrap(err, "searching elasticsearch configuration")
	}

	stack, err := GetStack(ctx, streamProcessor)
	if err != nil {
		return err
	}

	env := []corev1.EnvVar{
		Env("OPENSEARCH_URL", elasticSearchConfiguration.Spec.Endpoint()),
		Env("TOPIC_PREFIX", streamProcessor.Spec.Stack+"-"),
		Env("OPENSEARCH_INDEX", "stacks"),
		Env("STACK", streamProcessor.Spec.Stack),
	}
	if streamProcessor.Spec.Batching != nil {
		env = append(env,
			Env("OPENSEARCH_BATCHING_COUNT", fmt.Sprint(streamProcessor.Spec.Batching.Count)),
			Env("OPENSEARCH_BATCHING_PERIOD", streamProcessor.Spec.Batching.Period),
		)
	}
	configuration, err := GetByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, streamProcessor.Spec.Stack)
	if err != nil {
		return err
	}
	if configuration != nil {
		env = append(env, opentelemetryconfigurations.GetEnvVars(configuration, "gateway")...)
	}
	if brokerConfiguration.Spec.Kafka != nil {
		env = append(env, Env("KAFKA_ADDRESS", strings.Join(brokerConfiguration.Spec.Kafka.Brokers, ",")))
		if brokerConfiguration.Spec.Kafka.TLS {
			env = append(env, Env("KAFKA_TLS_ENABLED", "true"))
		}
		if brokerConfiguration.Spec.Kafka.SASL != nil {
			env = append(env,
				Env("KAFKA_SASL_USERNAME", brokerConfiguration.Spec.Kafka.SASL.Username),
				Env("KAFKA_SASL_PASSWORD", brokerConfiguration.Spec.Kafka.SASL.Password),
				Env("KAFKA_SASL_MECHANISM", brokerConfiguration.Spec.Kafka.SASL.Mechanism),
			)
		}
	}
	if brokerConfiguration.Spec.Nats != nil {
		env = append(env, Env("NATS_URL", brokerConfiguration.Spec.Nats.URL))
	}
	if elasticSearchConfiguration.Spec.BasicAuth != nil {
		env = append(env, Env("BASIC_AUTH_ENABLED", "true"))
		if elasticSearchConfiguration.Spec.BasicAuth.SecretName == "" {
			env = append(env,
				Env("BASIC_AUTH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.Username),
				Env("BASIC_AUTH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.Password),
			)
		} else {
			env = append(env,
				EnvFromSecret("BASIC_AUTH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "username"),
				EnvFromSecret("BASIC_AUTH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "password"),
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
	isOpenTelemetryEnabled, err := IsEnabledByLabel[*v1beta1.OpenTelemetryConfiguration](ctx, streamProcessor.Spec.Stack)
	if err != nil {
		return err
	}
	if isOpenTelemetryEnabled {
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
	if isOpenTelemetryEnabled {
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
		object.SetName("stream-processor-audit")
		if err := client.IgnoreNotFound(ctx.GetClient().Delete(ctx, object)); err != nil {
			return errors.Wrap(err, "deleting audit config map")
		}
	}

	configMaps := make([]*corev1.ConfigMap, 0)

	for _, x := range directories {
		data := make(map[string]string)

		CopyDir(x.fs, x.name, x.name, &data)

		configMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
			Namespace: streamProcessor.Spec.Stack,
			Name:      "stream-processor-" + x.name,
		},
			func(t *corev1.ConfigMap) {
				t.Data = data
			},
			WithController[*corev1.ConfigMap](ctx.GetScheme(), streamProcessor),
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
						Name: "stream-processor-" + x.name,
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

	streamList := &v1beta1.StreamList{}
	if err := ctx.GetClient().List(ctx, streamList, client.MatchingFields{
		"stack": streamProcessor.Spec.Stack,
	}); err != nil {
		return err
	}

	streams := streamList.Items
	sort.Slice(streams, func(i, j int) bool {
		return streams[i].Name < streams[j].Name
	})

	_, _, err = CreateOrUpdate(ctx, types.NamespacedName{
		Namespace: streamProcessor.Spec.Stack,
		Name:      "stream-processor",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), streamProcessor),
		deployments.WithMatchingLabels("stream-processor"),
		deployments.WithInitContainers(streamProcessor.Spec.InitContainers...),
		deployments.WithContainers(corev1.Container{
			Name:    "stream-processor",
			Image:   "public.ecr.aws/formance-internal/jeffail/benthos:v4.23.0-es",
			Env:     env,
			Command: cmd,
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
					Sources: Map(streams, func(stream v1beta1.Stream) corev1.VolumeProjection {
						return corev1.VolumeProjection{
							ConfigMap: &corev1.ConfigMapProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: fmt.Sprintf("stream-%s", stream.Name),
								},
								Items: []corev1.KeyToPath{{
									Key:  "stream.yaml",
									Path: stream.Name + ".yaml",
								}},
							},
						}
					}),
				},
			},
		})...),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Annotations = MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"config-hash": HashFromConfigMaps(configMaps...),
			})
		},
	)
	return err
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithWatchConfigurationObject(&v1beta1.BrokerConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.ElasticSearchConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.RegistriesConfiguration{}),
			WithWatchStack(),
			WithWatchDependency(&v1beta1.Stream{}),
			WithOwn(&corev1.ConfigMap{}),
			WithOwn(&appsv1.Deployment{}),
		),
	)
}
