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
	"strconv"

	"github.com/formancehq/operator/internal/resources/secretreferences"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	deployments "github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/httpapis"
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
	elasticSearchURI, err := settings.RequireURL(ctx, stack.Name, "elasticsearch.dsn")
	if err != nil {
		return err
	}

	secretReference, err := secretreferences.Sync(ctx, search, "elasticsearch", elasticSearchURI)
	if err != nil {
		return err
	}

	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, LowerCamelCaseName(ctx, search))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, GetDevEnvVars(stack, search)...)

	env = append(env,
		Env("OPEN_SEARCH_SERVICE", elasticSearchURI.Host),
		Env("OPEN_SEARCH_SCHEME", elasticSearchURI.Scheme),
		Env("ES_INDICES", "stacks"),
	)
	if secret := elasticSearchURI.Query().Get("secret"); elasticSearchURI.User != nil || secret != "" {
		if secret == "" {
			password, _ := elasticSearchURI.User.Password()
			env = append(env,
				Env("OPEN_SEARCH_USERNAME", elasticSearchURI.User.Username()),
				Env("OPEN_SEARCH_PASSWORD", password),
			)
		} else {
			secret := fmt.Sprintf("%s-%s", search.Name, secret)
			env = append(env,
				EnvFromSecret("OPEN_SEARCH_USERNAME", secret, "username"),
				EnvFromSecret("OPEN_SEARCH_PASSWORD", secret, "password"),
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

		batchingMap, err := settings.GetMapOrEmpty(ctx, stack.Name, "search.batching")
		if err != nil {
			return err
		}

		batching = &v1beta1.Batching{}
		if countString, ok := batchingMap["count"]; ok {
			count, err := strconv.ParseUint(countString, 10, 64)
			if err != nil {
				return err
			}
			batching.Count = int(count)
		}

		if period, ok := batchingMap["period"]; ok {
			batching.Period = period
		}
	}

	_, _, err = CreateOrUpdate[*v1beta1.Benthos](ctx, types.NamespacedName{
		Name: GetObjectName(stack.Name, "benthos"),
	},
		WithController[*v1beta1.Benthos](ctx.GetScheme(), search),
		func(t *v1beta1.Benthos) error {
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

	_, err = deployments.CreateOrUpdate(ctx, stack, search, "search",
		secretreferences.Annotate[*appsv1.Deployment](secretReference),
		deployments.WithMatchingLabels("search"),
		deployments.WithContainers(corev1.Container{
			Name:          "search",
			Image:         image,
			Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
			Env:           env,
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)

	if err := httpapis.Create(ctx, search); err != nil {
		return err
	}

	return err
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Search](),
			WithOwn[*v1beta1.Search](&v1beta1.SecretReference{}),
			WithOwn[*v1beta1.Search](&v1beta1.Benthos{}),
			WithOwn[*v1beta1.Search](&v1beta1.HTTPAPI{}),
			WithOwn[*v1beta1.Search](&appsv1.Deployment{}),
		),
	)
}
