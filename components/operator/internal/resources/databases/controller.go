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

package databases

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/secrets"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	databaseFinalizer = "finalize.databases.formance.com"
)

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database) error {

	databaseConfiguration, err := settings.FindDatabaseConfiguration(ctx, database)
	if err != nil {
		return err
	}

	if databaseConfiguration == nil {
		return fmt.Errorf("unable to find a database configuration")
	}

	if databaseConfiguration.CredentialsFromSecret != "" {
		secret, err := secrets.SyncOne(ctx, database, stack.Name, databaseConfiguration.CredentialsFromSecret)
		if err != nil {
			return err
		}
		databaseConfiguration.CredentialsFromSecret = secret.Name
	}

	switch {
	// TODO: We have multiple occurrences of this type of code, we need to factorize
	case !database.Status.Ready:
		// Some job fields are immutable (env vars for example)
		// So, if the configuration has changed, wee need to delete the job,
		// Then recreate a new one
		if database.Status.Configuration != nil {
			if !reflect.DeepEqual(
				database.Status.Configuration.DatabaseConfiguration,
				*databaseConfiguration,
			) {
				object := &batchv1.Job{}
				object.SetName(fmt.Sprintf("%s-create-database", database.Spec.Service))
				object.SetNamespace(database.Spec.Stack)
				if err := ctx.GetClient().Delete(ctx, object, &client.DeleteOptions{
					GracePeriodSeconds: pointer.For(int64(0)),
				}); client.IgnoreNotFound(err) != nil {
					return err
				}
			}
		}

		dbName := core.GetObjectName(database.Spec.Stack, database.Spec.Service)
		job, err := createDatabase(ctx, *databaseConfiguration, database, dbName)
		if err != nil {
			return err
		}

		database.Status.Configuration = &v1beta1.CreatedDatabase{
			DatabaseConfiguration: *databaseConfiguration,
			Database:              dbName,
		}

		if job.Status.Succeeded == 0 {
			return core.NewPendingError()
		}
	case !reflect.DeepEqual(database.Status.Configuration.DatabaseConfiguration,
		*databaseConfiguration):
		database.Status.OutOfSync = true
	}

	return nil
}

func Delete(ctx core.Context, database *v1beta1.Database) error {
	if database.Status.Configuration == nil {
		return nil
	}
	job, _, err := core.CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-drop-database", database.Spec.Service),
	},
		func(t *batchv1.Job) error {
			dropDBCommand := `psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME} -c "DROP DATABASE \"${POSTGRES_DATABASE}\""`
			if database.Status.Configuration.DisableSSLMode {
				dropDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:  "drop-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", dropDBCommand},
				Env: append(PostgresEnvVars(database.Status.Configuration.DatabaseConfiguration,
					database.Status.Configuration.Database),
					// psql use PGPASSWORD env var
					core.Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}

			return nil
		},
		core.WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	if err != nil {
		return err
	}

	if job.Status.Succeeded == 0 {
		return core.NewPendingError()
	}

	return nil
}

func createDatabase(ctx core.Context, databaseConfiguration v1beta1.DatabaseConfiguration,
	database *v1beta1.Database, dbName string) (*batchv1.Job, error) {

	job, _, err := core.CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-database", database.Spec.Service),
	},
		func(t *batchv1.Job) error {
			// PG does not support 'CREATE IF NOT EXISTS ' construct, emulate it with the above query
			createDBCommand := `echo SELECT \'CREATE DATABASE \"${POSTGRES_DATABASE}\"\' WHERE NOT EXISTS \(SELECT FROM pg_database WHERE datname = \'${POSTGRES_DATABASE}\'\)\\gexec | psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME}`
			if databaseConfiguration.DisableSSLMode {
				createDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:  "create-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", createDBCommand},
				Env: append(PostgresEnvVars(databaseConfiguration, dbName),
					// psql use PGPASSWORD env var
					core.Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}

			return nil
		},
		core.WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	return job, err
}

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn[*v1beta1.Database](&batchv1.Job{}),
			core.WithOwn[*v1beta1.Database](&v1.Secret{}),
			core.WithWatchSettings[*v1beta1.Database](),
			core.WithWatchSecrets[*v1beta1.Database](),
			core.WithFinalizer(databaseFinalizer, Delete),
		),
	)
}
