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

package formance_com

import (
	"fmt"
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/authclients"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/databases"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	. "github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/formancehq/operator/v2/internal/resources/topicqueries"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
)

// OrchestrationController reconciles a Orchestration object
type OrchestrationController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/finalizers,verbs=update

func (r *OrchestrationController) Reconcile(ctx Context, orchestration *v1beta1.Orchestration) error {

	stack, err := stacks.GetStack(ctx, orchestration.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, orchestration)
	if err != nil {
		return err
	}

	authClient, err := r.handleAuthClient(ctx, stack, orchestration)
	if err != nil {
		return err
	}

	if err := r.handleTopics(ctx, stack, orchestration); err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, orchestration, database, authClient); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, orchestration, "orchestration"); err != nil {
		return err
	}

	return nil
}

func (r *OrchestrationController) handleAuthClient(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration) (*v1beta1.AuthClient, error) {

	auth, err := stacks.GetIfEnabled[*v1beta1.Auth](ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, errors.New("requires auth service")
	}

	return authclients.Create(ctx, stack, orchestration, "orchestration",
		func(spec *v1beta1.AuthClientSpec) {
			spec.Scopes = []string{
				"ledger:read",
				"ledger:write",
				"payments:read",
				"payments:write",
				"wallets:read",
				"wallets:write",
			}
		})
}

func (r *OrchestrationController) handleTopics(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration) error {
	return ForEachEventPublisher(ctx, stack.Name, func(object client.Object) error {
		return topicqueries.Create(ctx, stack, strings.ToLower(object.GetObjectKind().GroupVersionKind().Kind), orchestration)
	})
}

func (r *OrchestrationController) createDeployment(ctx Context, stack *v1beta1.Stack,
	orchestration *v1beta1.Orchestration, database *v1beta1.Database, client *v1beta1.AuthClient) error {
	env, err := GetCommonServicesEnvVars(ctx, stack, "orchestration", orchestration.Spec)
	if err != nil {
		return err
	}
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)

	eventPublishers, err := ListEventPublishers(ctx, stack.Name)
	if err != nil {
		return err
	}

	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", orchestration.Spec.Temporal.Address),
		Env("TEMPORAL_NAMESPACE", orchestration.Spec.Temporal.Namespace),
		Env("WORKER", "true"),
		Env("TOPICS", strings.Join(Map(eventPublishers, func(from unstructured.Unstructured) string {
			return fmt.Sprintf("%s-%s", stack.Name, strings.ToLower(from.GetKind()))
		}), " ")),
	)

	env = append(env, authclients.GetEnvVars(client)...)

	if orchestration.Spec.Temporal.TLS.SecretName == "" {
		env = append(env,
			Env("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.Key),
			Env("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.CRT),
		)
	} else {
		env = append(env,
			EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.SecretName, "tls.key"),
			EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.SecretName, "tls.crt"),
		)
	}

	brokerEnvVars, err := brokerconfigurations.GetEnvVars(ctx, stack.Name, "orchestration")
	if err != nil && !errors.Is(err, stacks.ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	image, err := GetImage(ctx, stack, "orchestration", orchestration.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: orchestration.Spec.Stack,
		Name:      "orchestration",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), orchestration),
		deployments.WithMatchingLabels("orchestration"),
		deployments.WithContainers(corev1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Resources:     GetResourcesWithDefault(orchestration.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrchestrationController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(
				databases.Watch[*v1beta1.Orchestration](mgr, "orchestration")),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.Ledger{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.Payments{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.Wallets{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.RegistriesConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Orchestration](mgr)),
		).
		Owns(&v1beta1.TopicQuery{}).
		Owns(&v1beta1.AuthClient{}).
		For(&v1beta1.Orchestration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForOrchestration() *OrchestrationController {
	return &OrchestrationController{}
}
