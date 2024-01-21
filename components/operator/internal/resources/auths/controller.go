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

package auths

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/httpapis"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=auths,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=auths/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=auths/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, auth *v1beta1.Auth, version string) error {

	authClientList := make([]*v1beta1.AuthClient, 0)
	err := GetAllStackDependencies(ctx, auth.Spec.Stack, &authClientList)
	if err != nil {
		return err
	}

	configMap, err := createConfiguration(ctx, stack, authClientList)
	if err != nil {
		return errors.Wrap(err, "creating configuration")
	}

	database, err := databases.Create(ctx, auth)
	if err != nil {
		return errors.Wrap(err, "creating database")
	}

	if database.Status.Ready {
		_, err := createDeployment(ctx, stack, auth, database, configMap, version)
		if err != nil {
			return errors.Wrap(err, "creating deployment")
		}
	}

	if err := httpapis.Create(ctx, auth,
		httpapis.WithRules(httpapis.RuleUnsecured()),
		httpapis.WithServiceConfiguration(auth.Spec.Service)); err != nil {
		return errors.Wrap(err, "creating http api")
	}

	auth.Status.Clients = Map(authClientList, (*v1beta1.AuthClient).GetName)

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn(&appsv1.Deployment{}),
			WithOwn(&v1beta1.HTTPAPI{}),
			WithOwn(&v1beta1.Database{}),
			WithOwn(&corev1.ConfigMap{}),
			WithWatch(databases.Watch("auth", &v1beta1.Auth{})),
			WithWatchConfigurationObject(&v1beta1.RegistriesConfiguration{}),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchDependency(&v1beta1.AuthClient{}),
			WithWatchStack(),
		),
	)
}
