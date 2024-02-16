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

package payments

import (
	_ "embed"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	"net/http"

	"github.com/formancehq/operator/internal/resources/benthosstreams"
	"github.com/formancehq/search/benthos"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=payments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=payments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=payments/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, p *v1beta1.Payments, version string) error {

	database, err := databases.Create(ctx, stack, p)
	if err != nil {
		return err
	}

	if database.Status.Ready {
		image, err := registries.GetImage(ctx, stack, "payments", version)
		if err != nil {
			return err
		}

		if databases.GetSavedModuleVersion(database) != version {
			encryptionKey, err := getEncryptionKey(ctx, p)
			if err != nil {
				return err
			}

			if err := jobs.Handle(ctx, p, "migrate",
				databases.MigrateDatabaseContainer(image, database,
					func(m *databases.MigrationConfiguration) {
						m.AdditionalEnv = []corev1.EnvVar{
							Env("CONFIG_ENCRYPTION_KEY", encryptionKey),
						}
					},
				),
				jobs.WithServiceAccount(database.Status.URI.Query().Get("awsRole")),
			); err != nil {
				return err
			}

			if err := databases.SaveModuleVersion(ctx, database, version); err != nil {
				return errors.Wrap(err, "saving module version in database object")
			}
		}

		if semver.IsValid(version) && semver.Compare(version, "v1.0.0-alpha") < 0 {
			if err := createFullDeployment(ctx, stack, p, database, image); err != nil {
				return err
			}
		} else {
			if err := createReadDeployment(ctx, stack, p, database, image); err != nil {
				return err
			}

			if err := createConnectorsDeployment(ctx, stack, p, database, image); err != nil {
				return err
			}
			if err := createGateway(ctx, stack, p); err != nil {
				return err
			}
		}
	}

	search := &v1beta1.Search{}
	hasSearch, err := HasDependency(ctx, stack.Name, search)
	if err != nil {
		return err
	}
	if hasSearch {
		if err := benthosstreams.LoadFromFileSystem(ctx, benthos.Streams, p, "streams/payments/v0.0.0"); err != nil {
			return err
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.BenthosStream{}, client.MatchingLabels{
			"service": "payments",
		}); err != nil {
			return err
		}
	}

	if err := gatewayhttpapis.Create(ctx, p,
		gatewayhttpapis.WithHealthCheckEndpoint("_health"),
		gatewayhttpapis.WithRules(
			v1beta1.GatewayHTTPAPIRule{
				Path:    "/connectors/webhooks",
				Methods: []string{http.MethodPost},
				Secured: true,
			},
			gatewayhttpapis.RuleSecured(),
		)); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn[*v1beta1.Payments](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Payments](&corev1.Service{}),
			WithOwn[*v1beta1.Payments](&v1beta1.GatewayHTTPAPI{}),
			WithOwn[*v1beta1.Payments](&batchv1.Job{}),
			WithWatchSettings[*v1beta1.Payments](),
			WithWatchDependency[*v1beta1.Payments](&v1beta1.Search{}),
			databases.Watch[*v1beta1.Payments](),
			brokertopics.Watch[*v1beta1.Payments]("payments"),
		),
	)
}
