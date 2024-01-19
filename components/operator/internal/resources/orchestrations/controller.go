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

package orchestrations

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopicconsumers"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/httpapis"
	appsv1 "k8s.io/api/apps/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, o *v1beta1.Orchestration, version string) error {

	database, err := databases.Create(ctx, o)
	if err != nil {
		return err
	}

	authClient, err := createAuthClient(ctx, stack, o)
	if err != nil {
		return err
	}

	consumers, err := brokertopicconsumers.CreateOrUpdateOnAllServices(ctx, o)
	if err != nil {
		return err
	}

	if err := httpapis.Create(ctx, o,
		httpapis.WithServiceConfiguration(o.Spec.Service)); err != nil {
		return err
	}

	if database.Status.Ready && consumers.Ready() {
		if err := createDeployment(ctx, stack, o, database, authClient, consumers, version); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn(&v1beta1.BrokerTopicConsumer{}),
			WithOwn(&v1beta1.AuthClient{}),
			WithOwn(&appsv1.Deployment{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithWatchStack(),
			WithWatch(databases.Watch("orchestration", &v1beta1.Orchestration{})),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.TemporalConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.RegistriesConfiguration{}),
			WithWatchDependency(&v1beta1.Ledger{}),
			WithWatchDependency(&v1beta1.Auth{}),
			WithWatchDependency(&v1beta1.Payments{}),
			WithWatchDependency(&v1beta1.Wallets{}),
		),
	)
}
