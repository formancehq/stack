/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	common "github.com/formancehq/operator/v2/internal/common"
	databases "github.com/formancehq/operator/v2/internal/databases"
	"github.com/formancehq/operator/v2/internal/deployments"
	"github.com/formancehq/operator/v2/internal/httpapis"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	. "github.com/formancehq/operator/v2/internal/utils"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
)

// AuthController reconciles a Auth object
type AuthController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=auths,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=auths/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=auths/finalizers,verbs=update

func (r *AuthController) Reconcile(ctx reconcilers.Context, auth *v1beta1.Auth) error {

	stack, err := common.GetStack(ctx, auth.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, stack, "auth")
	if err != nil {
		return err
	}

	authClientsList := &v1beta1.AuthClientList{}
	if err := ctx.GetClient().List(ctx, authClientsList, client.MatchingFields{
		".spec.stack": stack.Name,
	}); err != nil {
		return err
	}

	configMap, err := r.createConfiguration(ctx, stack, authClientsList.Items)
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, auth, database, configMap); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, auth, "auth", httpapis.Secured()); err != nil {
		return err
	}

	auth.Status.Clients = Map(authClientsList.Items, func(from v1beta1.AuthClient) string {
		return from.Name
	})
	return nil
}

func (r *AuthController) createConfiguration(ctx reconcilers.Context, stack *v1beta1.Stack, items []v1beta1.AuthClient) (*corev1.ConfigMap, error) {

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	yamlData, err := yaml.Marshal(struct {
		Clients any `yaml:"clients"`
	}{
		Clients: Map(items, func(from v1beta1.AuthClient) any {
			return from.Spec
		}),
	})
	if err != nil {
		return nil, err
	}

	cm, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "auth-configuration",
	}, func(t *corev1.ConfigMap) {
		t.Data = map[string]string{
			"config.yaml": string(yamlData),
		}
	})
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func (r *AuthController) createDeployment(ctx reconcilers.Context, stack *v1beta1.Stack, auth *v1beta1.Auth, database *v1beta1.Database, configMap *corev1.ConfigMap) error {

	env, err := common.GetCommonServicesEnvVars(ctx, stack, "auth", auth.Spec)
	if err != nil {
		return err
	}

	env = append(env,
		databases.PostgresEnvVars(
			database.Status.Configuration.DatabaseConfigurationSpec,
			common.GetObjectName(stack.Name, "auth"),
		)...,
	)
	env = append(env,
		Env("CONFIG", "/config/config.yaml"),
		Env("BASE_URL", "$(STACK_PUBLIC_URL)/api/auth"),
	)
	if auth.Spec.SigningKey != "" && auth.Spec.SigningKeyFromSecret != nil {
		return fmt.Errorf("cannot specify signing key using both .spec.signingKey and .spec.signingKeyFromSecret fields")
	}
	if auth.Spec.SigningKey != "" {
		env = append(env, Env("SIGNING_KEY", auth.Spec.SigningKey))
	}
	if auth.Spec.SigningKeyFromSecret != nil {
		env = append(env, corev1.EnvVar{
			Name: "SIGNING_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: auth.Spec.SigningKeyFromSecret,
			},
		})
	}
	if auth.Spec.DelegatedOIDCServer != nil {
		env = append(env,
			Env("DELEGATED_CLIENT_SECRET", auth.Spec.DelegatedOIDCServer.ClientSecret),
			Env("DELEGATED_CLIENT_ID", auth.Spec.DelegatedOIDCServer.ClientID),
			Env("DELEGATED_ISSUER", auth.Spec.DelegatedOIDCServer.Issuer),
		)
	}
	if stack.Spec.Dev || auth.Spec.Dev {
		env = append(env, Env("CAOS_OIDC_DEV", "1"))
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx,
		common.GetNamespacedResourceName(stack.Name, "auth"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Annotations = MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"config-hash": HashFromConfigMap(configMap),
			})
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:      "auth",
				Args:      []string{"serve"},
				Env:       env,
				Image:     common.GetImage("auth", common.GetVersion(stack, auth.Spec.Version)),
				Resources: GetResourcesWithDefault(auth.Spec.ResourceProperties, ResourceSizeSmall()),
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "config",
					ReadOnly:  true,
					MountPath: "/config",
				}},
				Ports: []corev1.ContainerPort{common.StandardHTTPPort()},
			}}
			t.Spec.Template.Spec.Volumes = []corev1.Volume{{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: configMap.Name,
						},
					},
				},
			}}
		},
		deployments.WithMatchingLabels("auth"),
		WithController[*appsv1.Deployment](ctx.GetScheme(), auth),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *AuthController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Auth{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Auth).Spec.Stack}
	}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Auth{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&v1beta1.Database{}).
		Watches(
			&v1beta1.AuthClient{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				authClient := object.(*v1beta1.AuthClient)

				list := &v1beta1.AuthList{}
				if err := mgr.GetClient().List(ctx, list, client.MatchingFields{
					".spec.stack": authClient.Spec.Stack,
				}); err != nil {
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.Auth])...,
				)
			}),
		), nil
}

func ForAuth() *AuthController {
	return &AuthController{}
}
