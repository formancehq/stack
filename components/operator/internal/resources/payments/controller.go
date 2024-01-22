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
	"net/http"

	"github.com/formancehq/operator/internal/resources/streams"
	"github.com/formancehq/search/benthos"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/httpapis"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=payments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=payments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=payments/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, p *v1beta1.Payments, version string) error {

	database, err := databases.Create(ctx, p)
	if err != nil {
		return err
	}

	if database.Status.Ready {
		if semver.IsValid(version) && semver.Compare(version, "v1.0.0-alpha") < 0 {
			if err := createFullDeployment(ctx, stack, p, database, version); err != nil {
				return err
			}
		} else {
			if err := createReadDeployment(ctx, stack, p, database, version); err != nil {
				return err
			}

			if err := createConnectorsDeployment(ctx, stack, p, database, version); err != nil {
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
		if err := streams.LoadFromFileSystem(ctx, benthos.Streams, p.Spec.Stack, "streams/payments/v0.0.0",
			WithController[*v1beta1.Stream](ctx.GetScheme(), p),
			WithLabels[*v1beta1.Stream](map[string]string{
				"service": "payments",
			}),
		); err != nil {
			return err
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.Stream{}, client.MatchingLabels{
			"service": "payments",
		}); err != nil {
			return err
		}
	}

	if err := httpapis.Create(ctx, p,
		httpapis.WithRules(
			v1beta1.HTTPAPIRule{
				Path:    "/connectors/webhooks",
				Methods: []string{http.MethodPost},
				Secured: true,
			},
			httpapis.RuleSecured(),
		), httpapis.WithServiceConfiguration(p.Spec.Service)); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn(&appsv1.Deployment{}),
			WithOwn(&corev1.Service{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithWatchStack(),
			WithWatch(databases.Watch("payments", &v1beta1.Payments{})),
			WithWatch[*v1beta1.BrokerTopic](brokertopics.Watch[*v1beta1.Payments]("payments")),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchDependency(&v1beta1.Search{}),
		),
	)
}
