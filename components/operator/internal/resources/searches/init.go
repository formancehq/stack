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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconsumers"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	"github.com/formancehq/operator/internal/resources/gateways"
	. "github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"github.com/formancehq/operator/internal/resources/settings"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, search *v1beta1.Search, version string) error {
	elasticSearchURI, err := settings.RequireURL(ctx, stack.Name, "elasticsearch", "dsn")
	if err != nil {
		return err
	}

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	awsIAMEnabled := serviceAccountName != ""

	var elasticSearchSecretResourceRef *v1beta1.ResourceReference
	if secret := elasticSearchURI.Query().Get("secret"); !awsIAMEnabled && secret != "" {
		elasticSearchSecretResourceRef, err = resourcereferences.Create(ctx, search, "elasticsearch", secret, &corev1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, search, "elasticsearch")
	}
	if err != nil {
		return err
	}

	env := make([]corev1.EnvVar, 0)
	if awsIAMEnabled {
		env = append(env, Env("AWS_IAM_ENABLED", "true"))
	}

	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, LowerCamelCaseKind(ctx, search))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, GetDevEnvVars(stack, search)...)

	gatewayEnvVars, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnvVars...)

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

	if err := createConsumers(ctx, search); err != nil {
		return err
	}

	batching := search.Spec.Batching
	if batching == nil {

		batchingMap, err := settings.GetMapOrEmpty(ctx, stack.Name, "search", "batching")
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

	annotations := map[string]string{}
	if elasticSearchSecretResourceRef != nil {
		annotations["elasticsearch-secret-hash"] = elasticSearchSecretResourceRef.Status.Hash
	}

	err = applications.
		New(search, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "search",
			},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: annotations,
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: []corev1.Container{{
							Name:          "search",
							Image:         image,
							Ports:         []corev1.ContainerPort{applications.StandardHTTPPort()},
							Env:           env,
							LivenessProbe: applications.DefaultLiveness("http"),
						}},
					},
				},
			},
		}).
		IsEE().
		Install(ctx)
	if err != nil {
		return err
	}

	if err := gatewayhttpapis.Create(ctx, search, gatewayhttpapis.WithHealthCheckEndpoint("_healthcheck")); err != nil {
		return err
	}

	if !search.Status.TopicCleaned {
		if err := cleanConsumers(ctx, search); err != nil {
			return fmt.Errorf("failed to clean consumers for search: %w", err)
		}
		search.Status.TopicCleaned = true
	}

	return err
}

func createConsumers(ctx Context, search *v1beta1.Search) error {
	for _, o := range []v1beta1.Module{
		&v1beta1.Payments{},
		&v1beta1.Ledger{},
		&v1beta1.Gateway{},
	} {
		if ok, err := HasDependency(ctx, search.Spec.Stack, o); err != nil {
			return err
		} else if ok {
			consumer, err := brokerconsumers.Create(ctx, search, LowerCamelCaseKind(ctx, o), LowerCamelCaseKind(ctx, o))
			if err != nil {
				return err
			}
			if !consumer.Status.Ready {
				return NewPendingError().WithMessage("waiting for consumer %s to be ready", consumer.Name)
			}
		}
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Search](),
			WithOwn[*v1beta1.Search](&v1beta1.BrokerConsumer{}),
			WithOwn[*v1beta1.Search](&v1beta1.ResourceReference{}),
			WithOwn[*v1beta1.Search](&v1beta1.Benthos{}),
			WithOwn[*v1beta1.Search](&v1beta1.GatewayHTTPAPI{}),
			WithOwn[*v1beta1.Search](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Search](&v1.Job{}),
		),
	)
}
