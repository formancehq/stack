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
	. "github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, search *v1beta1.Search, version string) error {
	elasticSearchConfiguration, err := settings.FindElasticSearchConfiguration(ctx, stack)
	if err != nil {
		return err
	}

	if elasticSearchConfiguration == nil {
		return errors.New("missing elastic search configuration")
	}

	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVarsIfEnabled(ctx, stack, GetModuleName(ctx, search))
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

		batchingCount, err := settings.GetIntOrDefault(ctx, stack.Name, 0, "search.batching.count")
		if err != nil {
			return err
		}

		batchingPeriod, err := settings.GetStringOrEmpty(ctx, stack.Name, "search.batching.period")
		if err != nil {
			return err
		}

		batching = &v1beta1.Batching{
			Count:  batchingCount,
			Period: batchingPeriod,
		}
	}

	_, _, err = CreateOrUpdate[*v1beta1.StreamProcessor](ctx, types.NamespacedName{
		Name: GetObjectName(stack.Name, "stream-processor"),
	},
		WithController[*v1beta1.StreamProcessor](ctx.GetScheme(), search),
		func(t *v1beta1.StreamProcessor) error {
			t.Spec.Stack = stack.Name
			t.Spec.Batching = batching
			t.Spec.DevProperties = search.Spec.DevProperties
			t.Spec.InitContainers = []corev1.Container{{
				Name:  "init-mapping",
				Image: image,
				Args:  []string{"init-mapping"},
				Env:   env,
			}}

			return nil
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
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)

	if err := httpapis.Create(ctx, search,
		httpapis.WithServiceConfiguration(search.Spec.Service)); err != nil {
		return err
	}

	return err
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchStack(),
			WithWatchConfigurationObject(&v1beta1.Settings{}),
			WithOwn(&v1beta1.StreamProcessor{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithOwn(&appsv1.Deployment{}),
		),
	)
}
