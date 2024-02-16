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

package webhooks

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopicconsumers"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=webhooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, version string) error {
	database, err := databases.Create(ctx, stack, webhooks)
	if err != nil {
		return err
	}

	consumers, err := brokertopicconsumers.CreateOrUpdateOnAllServices(ctx, webhooks)
	if err != nil {
		return err
	}

	if err := gatewayhttpapis.Create(ctx, webhooks); err != nil {
		return err
	}

	if database.Status.Ready {
		image, err := registries.GetImage(ctx, stack, "webhooks", version)
		if err != nil {
			return errors.Wrap(err, "resolving image")
		}

		if IsGreaterOrEqual(version, "v2.0.0-rc.5") && databases.GetSavedModuleVersion(database) != version {
			if err := jobs.Handle(ctx, webhooks, "migrate",
				databases.MigrateDatabaseContainer(image, database),
				jobs.WithServiceAccount(database.Status.URI.Query().Get("awsRole")),
			); err != nil {
				return err
			}
			if err := databases.SaveModuleVersion(ctx, database, version); err != nil {
				return errors.Wrap(err, "saving module version in database object")
			}
		}

		if consumers.Ready() {
			if IsGreaterOrEqual(version, "v0.7.1") {
				if err := createSingleDeployment(ctx, stack, webhooks, database, consumers, version); err != nil {
					return err
				}
			} else {
				if err := createDualDeployment(ctx, stack, webhooks, database, consumers, version); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn[*v1beta1.Webhooks](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Webhooks](&v1beta1.GatewayHTTPAPI{}),
			WithOwn[*v1beta1.Webhooks](&batchv1.Job{}),
			WithWatchSettings[*v1beta1.Webhooks](),
			WithWatchDependency[*v1beta1.Webhooks](&v1beta1.Ledger{}),
			WithWatchDependency[*v1beta1.Webhooks](&v1beta1.Payments{}),
			databases.Watch[*v1beta1.Webhooks](),
		),
	)
}
