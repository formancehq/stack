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

package ledgers

import (
	_ "embed"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/httpapis"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/streams"
	"github.com/formancehq/search/benthos"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//+kubebuilder:rbac:groups=formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, version string) error {

	database, err := databases.Create(ctx, ledger)
	if err != nil {
		return err
	}

	image, err := registries.GetImage(ctx, stack, "ledger", version)
	if err != nil {
		return err
	}

	if err := httpapis.Create(ctx, ledger,
		httpapis.WithServiceConfiguration(ledger.Spec.Service)); err != nil {
		return err
	}

	isV2 := false
	if !semver.IsValid(version) || semver.Compare(version, "v2.0.0-alpha") > 0 {
		isV2 = true
	}

	search := &v1beta1.Search{}
	hasSearch, err := HasDependency(ctx, stack.Name, search)
	if err != nil {
		return err
	}
	if hasSearch {
		streamsVersion := "v1.0.0"
		if isV2 {
			streamsVersion = "v2.0.0"
		}
		if err := streams.LoadFromFileSystem(ctx, benthos.Streams, ledger.Spec.Stack, "streams/ledger/"+streamsVersion,
			WithController[*v1beta1.Stream](ctx.GetScheme(), ledger),
			WithLabels[*v1beta1.Stream](map[string]string{
				"service": "ledger",
			}),
		); err != nil {
			return err
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.Stream{}, client.MatchingLabels{
			"service": "ledger",
		}); err != nil {
			return err
		}
	}

	if database.Status.Ready {

		actualVersion, err := ActualVersion(ctx, ledger)
		if err != nil {
			return err
		}

		actualVersionIsV1 := false
		if !semver.IsValid(actualVersion) || semver.Compare(actualVersion, "v2.0.0-alpha") < 0 {
			actualVersionIsV1 = true
		}

		if actualVersionIsV1 && isV2 {
			if err := migrateToLedgerV2(ctx, stack, ledger, database, image); err != nil {
				return err
			}
		}

		err = installLedger(ctx, stack, ledger, database, image, isV2)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn(&appsv1.Deployment{}),
			WithOwn(&batchv1.Job{}),
			WithOwn(&corev1.Service{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithOwn(&v1beta1.Database{}),
			WithWatchStack(),
			WithWatch[*v1beta1.BrokerTopic](brokertopics.Watch[*v1beta1.Ledger]("ledger")),
			WithWatch(databases.Watch("ledger", &v1beta1.Ledger{})),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchDependency(&v1beta1.Search{}),
		),
	)
}
