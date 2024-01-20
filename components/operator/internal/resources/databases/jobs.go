package databases

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	v12 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createJob(ctx core.Context, databaseConfiguration v1beta1.DatabaseConfiguration,
	database *v1beta1.Database, dbName string) (*v12.Job, error) {

	job, _, err := core.CreateOrUpdate[*v12.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-database", database.Spec.Service),
	},
		func(t *v12.Job) {
			// PG does not support 'CREATE IF NOT EXISTS ' construct, emulate it with the above query
			createDBCommand := `echo SELECT \'CREATE DATABASE \"${POSTGRES_DATABASE}\"\' WHERE NOT EXISTS \(SELECT FROM pg_database WHERE datname = \'${POSTGRES_DATABASE}\'\)\\gexec | psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME}`
			if databaseConfiguration.Spec.DisableSSLMode {
				createDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:  "create-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", createDBCommand},
				Env: append(PostgresEnvVars(databaseConfiguration.Spec, dbName),
					// psql use PGPASSWORD env var
					core.Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}
		},
		core.WithController[*v12.Job](ctx.GetScheme(), database),
	)
	return job, err
}

func deleteJob(ctx core.Context, database *v1beta1.Database) (*v12.Job, error) {
	job, _, err := core.CreateOrUpdate[*v12.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-drop-database", database.Spec.Service),
	},
		func(t *v12.Job) {
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
				Env: append(PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec,
					database.Status.Configuration.Database),
					// psql use PGPASSWORD env var
					core.Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
				),
			}}
		},
		core.WithController[*v12.Job](ctx.GetScheme(), database),
	)
	return job, err
}