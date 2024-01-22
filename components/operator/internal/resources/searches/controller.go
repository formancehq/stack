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

package searches

import (
	"fmt"
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	deployments "github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/httpapis"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	. "github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, search *v1beta1.Search, version string) error {
	elasticSearchConfiguration, err := FindElasticSearchConfiguration(ctx, stack)
	if err != nil {
		return err
	}

	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, GetModuleName(ctx, search))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, GetDevEnvVars(stack, search)...)

	env = append(env,
		Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d", elasticSearchConfiguration.Host, elasticSearchConfiguration.Port)),
		Env("OPEN_SEARCH_SCHEME", elasticSearchConfiguration.Scheme),
		Env("ES_INDICES", "stacks"),
	)
	if elasticSearchConfiguration.BasicAuth != nil {
		if elasticSearchConfiguration.BasicAuth.SecretName == "" {
			env = append(env,
				Env("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.BasicAuth.Username),
				Env("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.BasicAuth.Password),
			)
		} else {
			env = append(env,
				EnvFromSecret("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.BasicAuth.SecretName, "username"),
				EnvFromSecret("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.BasicAuth.SecretName, "password"),
			)
		}
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "search", search.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := GetImage(ctx, stack, "search", version)
	if err != nil {
		return err
	}

	batching := search.Spec.Batching
	if batching == nil {
		batchingConfiguration, err := GetConfigurationObject[*v1beta1.SearchBatchingConfiguration](ctx, search.Spec.Stack)
		if err != nil {
			return err
		}
		if batchingConfiguration != nil {
			batching = &batchingConfiguration.Spec.Batching
		}
	}

	_, _, err = CreateOrUpdate[*v1beta1.StreamProcessor](ctx, types.NamespacedName{
		Name: GetObjectName(stack.Name, "stream-processor"),
	},
		WithController[*v1beta1.StreamProcessor](ctx.GetScheme(), search),
		func(t *v1beta1.StreamProcessor) {
			t.Spec.Stack = stack.Name
			t.Spec.Batching = batching
			t.Spec.DevProperties = search.Spec.DevProperties
			if streamProcessor := search.Spec.StreamProcessor; streamProcessor != nil {
				t.Spec.ResourceProperties = search.Spec.StreamProcessor.ResourceRequirements
			}
			t.Spec.InitContainers = []corev1.Container{{
				Name:  "init-mapping",
				Image: image,
				Args:  []string{"init-mapping"},
				Env:   env,
			}}
		},
	)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, search, "search",
		deployments.WithMatchingLabels("search"),
		deployments.WithContainers(corev1.Container{
			Name:          "search",
			Image:         image,
			Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
			Env:           env,
			Resources:     GetResourcesRequirementsWithDefault(search.Spec.ResourceRequirements, ResourceSizeSmall()),
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)

	if err := httpapis.Create(ctx, search,
		httpapis.WithServiceConfiguration(search.Spec.Service)); err != nil {
		return err
	}

	return err
}

func FindElasticSearchConfiguration(ctx Context, stack *v1beta1.Stack) (*v1beta1.ElasticSearchConfiguration, error) {
	elasticSearchHost, err := settings.RequireString(ctx, stack.Name, "elasticsearch.host")
	if err != nil {
		return nil, err
	}

	elasticSearchScheme, err := settings.GetStringOrDefault(ctx, stack.Name, "https", "elasticsearch.scheme")
	if err != nil {
		return nil, err
	}

	elasticSearchPort, err := settings.GetUInt16OrDefault(ctx, stack.Name, 9200, "elasticsearch.port")
	if err != nil {
		return nil, err
	}

	elasticSearchTLSEnabled, err := settings.GetBoolOrFalse(ctx, stack.Name, "elasticsearch.tls.enabled")
	if err != nil {
		return nil, err
	}

	elasticSearchTLSSkipCertVerify, err := settings.GetBoolOrFalse(ctx, stack.Name, "elasticsearch.tls.skip-cert-verify")
	if err != nil {
		return nil, err
	}

	var basicAuth *v1beta1.ElasticSearchBasicAuthConfig
	basicAuthEnabled, err := settings.GetBoolOrFalse(ctx, stack.Name, "elasticsearch.basic-auth.enabled")
	if basicAuthEnabled {
		elasticSearchBasicAuthUsername, err := settings.GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.username")
		if err != nil {
			return nil, err
		}

		elasticSearchBasicAuthPassword, err := settings.GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.password")
		if err != nil {
			return nil, err
		}

		elasticSearchBasicAuthSecret, err := settings.GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.secret")
		if err != nil {
			return nil, err
		}

		basicAuth = &v1beta1.ElasticSearchBasicAuthConfig{
			Username:   elasticSearchBasicAuthUsername,
			Password:   elasticSearchBasicAuthPassword,
			SecretName: elasticSearchBasicAuthSecret,
		}
	}

	return &v1beta1.ElasticSearchConfiguration{
		Scheme: elasticSearchScheme,
		Host:   elasticSearchHost,
		Port:   elasticSearchPort,
		TLS: v1beta1.ElasticSearchTLSConfig{
			Enabled:        elasticSearchTLSEnabled,
			SkipCertVerify: elasticSearchTLSSkipCertVerify,
		},
		BasicAuth: basicAuth,
	}, nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchStack(),
			WithWatchConfigurationObject(&v1beta1.Settings{}),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.SearchBatchingConfiguration{}),
			WithOwn(&v1beta1.StreamProcessor{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithOwn(&appsv1.Deployment{}),
		),
	)
}
