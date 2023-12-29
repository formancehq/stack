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
	. "github.com/formancehq/operator/v2/internal/common"
	. "github.com/formancehq/operator/v2/internal/databases"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	. "github.com/formancehq/operator/v2/internal/utils"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/formancehq/operator/v2/api/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DatabaseController reconciles a Database object
type DatabaseController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func (r *DatabaseController) Reconcile(ctx reconcilers.Context, database *v1beta1.Database) error {

	serviceSelectorRequirement, err := labels.NewRequirement("formance.com/service", selection.In, []string{"any", database.Spec.Service})
	if err != nil {
		return err
	}

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", database.Spec.Stack})
	if err != nil {
		return err
	}

	selector := labels.NewSelector().Add(*serviceSelectorRequirement, *stackSelectorRequirement)

	databaseConfigurationList := &v1beta1.DatabaseConfigurationList{}
	if err := ctx.GetClient().List(ctx, databaseConfigurationList, &client.ListOptions{
		LabelSelector: selector,
	}); err != nil {
		return err
	}

	switch len(databaseConfigurationList.Items) {
	case 0:
		database.Status.Error = "unable to find a database configuration"
		database.Status.Ready = false
	case 1:
		if database.Status.BoundTo == "" || !database.Status.Ready {
			databaseConfiguration := databaseConfigurationList.Items[0]
			dbName := GetObjectName(database.Spec.Stack, database.Spec.Service)

			job, err := r.createJob(ctx, databaseConfiguration, database, dbName)
			if err != nil {
				return err
			}

			database.Status.Ready = job.Status.Succeeded > 0
			database.Status.Error = ""
			database.Status.Configuration = &v1beta1.CreatedDatabase{
				DatabaseConfigurationSpec: databaseConfiguration.Spec,
				Database:                  dbName,
			}
			database.Status.BoundTo = databaseConfiguration.Name
		} else if !reflect.DeepEqual(database.Status.Configuration.DatabaseConfigurationSpec, databaseConfigurationList.Items[0].Spec) {
			database.Status.OutOfSync = true
		}
	default:
		database.Status.Error = "multiple database configuration object found"
		database.Status.Ready = false
	}

	return nil
}

func (r *DatabaseController) createJob(ctx reconcilers.Context, databaseConfiguration v1beta1.DatabaseConfiguration,
	database *v1beta1.Database, dbName string) (*batchv1.Job, error) {

	job, _, err := CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-database", database.Spec.Service),
	},
		func(t *batchv1.Job) {
			// PG does not support 'CREATE IF NOT EXISTS ' construct, emulate it with the above query
			createDBCommand := `echo SELECT \'CREATE DATABASE \"${POSTGRES_DATABASE}\"\' WHERE NOT EXISTS \(SELECT FROM pg_database WHERE datname = \'${POSTGRES_DATABASE}\'\)\\gexec | psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME}`
			if databaseConfiguration.Spec.DisableSSLMode {
				createDBCommand += ` "sslmode=disable"`
			}

			t.Spec.Template.Spec = corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyOnFailure,
				Containers: []corev1.Container{{
					Name:  "create-database",
					Image: "postgres:15-alpine",
					Args:  []string{"sh", "-c", createDBCommand},
					Env: append(PostgresEnvVars(databaseConfiguration.Spec, dbName),
						// psql use PGPASSWORD env var
						Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
					),
				}},
			}
		},
		WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	return job, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Database{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.DatabaseConfiguration{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				list := &v1beta1.DatabaseList{}
				if err := mgr.GetClient().List(ctx, list); err != nil {
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.Database])...,
				)
			}),
		).
		Owns(&batchv1.Job{}), nil
}

func ForDatabase() *DatabaseController {
	return &DatabaseController{}
}
