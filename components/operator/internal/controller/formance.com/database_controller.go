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

package formance_com

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/stacks"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	databaseFinalizer = "finalize.databases.formance.com"
)

// DatabaseController reconciles a Database object
type DatabaseController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func (r *DatabaseController) Reconcile(ctx Context, database *v1beta1.Database) error {

	serviceSelectorRequirement, err := labels.NewRequirement(ServiceLabel, selection.In, []string{"any", database.Spec.Service})
	if err != nil {
		return err
	}

	stackSelectorRequirement, err := labels.NewRequirement(StackLabel, selection.In, []string{"any", database.Spec.Stack})
	if err != nil {
		return err
	}

	selector := labels.NewSelector().Add(*serviceSelectorRequirement, *stackSelectorRequirement)

	if !database.DeletionTimestamp.IsZero() {
		if database.Status.Configuration != nil {
			job, err := r.createDeleteJob(ctx, database)
			if err != nil {
				return err
			}

			if job.Status.Succeeded == 0 {
				return ErrPending
			}
		}

		patch := client.MergeFrom(database.DeepCopy())
		if updated := controllerutil.RemoveFinalizer(database, databaseFinalizer); updated {
			if err := ctx.GetClient().Patch(ctx, database, patch); err != nil {
				return errors.Wrap(err, "removing finalizer")
			}
		}

		return ErrDeleted
	}

	databaseConfigurationList := &v1beta1.DatabaseConfigurationList{}
	if err := ctx.GetClient().List(ctx, databaseConfigurationList, &client.ListOptions{
		LabelSelector: selector,
	}); err != nil {
		return err
	}

	switch len(databaseConfigurationList.Items) {
	case 0:
		database.Status.Error = "unable to find a database configuration"
		database.Status.OutOfSync = true
	case 1:
		switch {
		case database.Status.BoundTo == "" || !database.Status.Ready:
			databaseConfiguration := databaseConfigurationList.Items[0]

			// Some job fields are immutable (env vars for example)
			// So, if the configuration has changed, wee need to delete the job,
			// Then recreate a new one
			if database.Status.Configuration != nil {
				if !reflect.DeepEqual(
					database.Status.Configuration.DatabaseConfigurationSpec,
					databaseConfiguration.Spec,
				) {
					object := &batchv1.Job{}
					object.SetNamespace(database.Spec.Stack)
					object.SetName(fmt.Sprintf("%s-create-database", database.Spec.Service))
					if err := ctx.GetClient().Delete(ctx, object); client.IgnoreNotFound(err) != nil {
						return err
					}
				}
			}

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

			patch := client.MergeFrom(database.DeepCopy())
			if updated := controllerutil.AddFinalizer(database, databaseFinalizer); updated {
				if err := ctx.GetClient().Patch(ctx, database, patch); err != nil {
					return err
				}
			}
		case !reflect.DeepEqual(database.Status.Configuration.DatabaseConfigurationSpec,
			databaseConfigurationList.Items[0].Spec):
			database.Status.OutOfSync = true
		}
	default:
		database.Status.Error = "multiple database configuration object found"
		database.Status.OutOfSync = true
	}

	return nil
}

func (r *DatabaseController) createJob(ctx Context, databaseConfiguration v1beta1.DatabaseConfiguration,
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

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:  "create-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", createDBCommand},
				Env: append(databases.PostgresEnvVars(databaseConfiguration.Spec, dbName),
					// psql use PGPASSWORD env var
					Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}
		},
		WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	return job, err
}

func (c DatabaseController) createDeleteJob(ctx Context, database *v1beta1.Database) (*batchv1.Job, error) {
	job, _, err := CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-drop-database", database.Spec.Service),
	},
		func(t *batchv1.Job) {
			dropDBCommand := `psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME} -c "DROP DATABASE \"${POSTGRES_DATABASE}\""`
			if database.Status.Configuration.DisableSSLMode {
				dropDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:  "drop-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", dropDBCommand},
				Env: append(databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec,
					database.Status.Configuration.Database),
					// psql use PGPASSWORD env var
					Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}
		},
		WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	return job, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Database{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.DatabaseConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Database](mgr)),
		).
		Watches(
			&v1beta1.RegistriesConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Database](mgr)),
		).
		Owns(&batchv1.Job{}), nil
}

func ForDatabase() *DatabaseController {
	return &DatabaseController{}
}
