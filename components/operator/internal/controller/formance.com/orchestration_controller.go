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
	appsv1 "k8s.io/api/apps/v1"
	"strings"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/brokertopicconsumers"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/httpapis"
	. "github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/stacks"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

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

	consumers, err := brokertopicconsumers.CreateOrUpdateOnAllServices(ctx, orchestration)
	if err != nil {
		return err
	}

	if err := httpapis.Create(ctx, orchestration,
		httpapis.WithServiceConfiguration(orchestration.Spec.Service)); err != nil {
		return err
	}

	if database.Status.Ready && consumers.Ready() {
		if err := r.createDeployment(ctx, stack, orchestration, database, authClient, consumers); err != nil {
			return err
		}
	}

	return nil
}

func (r *OrchestrationController) handleAuthClient(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration) (*v1beta1.AuthClient, error) {

	hasAuth, err := stacks.HasDependency[*v1beta1.Auth](ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return nil, nil
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

func (r *OrchestrationController) createDeployment(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration, database *v1beta1.Database, client *v1beta1.AuthClient, consumers []*v1beta1.BrokerTopicConsumer) error {
	env, err := GetCommonServicesEnvVars(ctx, stack, orchestration)
	if err != nil {
		return err
	}
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)

	temporalConfiguration, err := stacks.RequireLabelledConfig[*v1beta1.TemporalConfiguration](ctx, stack.Name)
	if err != nil {
		return err
	}

	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", temporalConfiguration.Spec.Address),
		Env("TEMPORAL_NAMESPACE", temporalConfiguration.Spec.Namespace),
		Env("WORKER", "true"),
		Env("TOPICS", strings.Join(Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
			return fmt.Sprintf("%s-%s", stack.Name, from.Spec.Service)
		}), " ")),
	)

	authEnvVars, err := auths.EnvVars(ctx, stack, "orchestration", orchestration.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	if client != nil {
		env = append(env, authclients.GetEnvVars(client)...)
	}

	if temporalConfiguration.Spec.TLS.SecretName == "" {
		env = append(env,
			Env("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.Spec.TLS.Key),
			Env("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.Spec.TLS.CRT),
		)
	} else {
		env = append(env,
			EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.Spec.TLS.SecretName, "tls.key"),
			EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.Spec.TLS.SecretName, "tls.crt"),
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

	_, err = deployments.CreateOrUpdate(ctx, orchestration, "orchestration",
		deployments.WithMatchingLabels("orchestration"),
		deployments.WithContainers(corev1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Resources:     GetResourcesRequirementsWithDefault(orchestration.Spec.ResourceRequirements, ResourceSizeSmall()),
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
			&v1beta1.TemporalConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.Ledger{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Orchestration](mgr)),
		).
		Watches(
			&v1beta1.Auth{},
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
		Owns(&v1beta1.BrokerTopicConsumer{}).
		Owns(&v1beta1.AuthClient{}).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		For(&v1beta1.Orchestration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForOrchestration() *OrchestrationController {
	return &OrchestrationController{}
}
