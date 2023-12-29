package databases

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	common "github.com/formancehq/operator/v2/internal/core"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func Create(ctx common.Context, stack *v1beta1.Stack, serviceName string) (*v1beta1.Database, error) {
	database, _, err := common.CreateOrUpdate[*v1beta1.Database](ctx, types.NamespacedName{
		Name: common.GetObjectName(stack.Name, serviceName),
	}, func(t *v1beta1.Database) {
		t.Spec.Stack = stack.Name
		t.Spec.Service = serviceName
	})
	if !database.Status.Ready {
		return nil, common.ErrPending
	}
	return database, err
}

type MigrationConfiguration struct {
	Command       []string
	AdditionalEnv []v1.EnvVar
}

func MigrateDatabaseContainer(serviceName string, database v1beta1.DatabaseConfigurationSpec, databaseName, version string, options ...func(m *MigrationConfiguration)) v1.Container {
	m := &MigrationConfiguration{}
	for _, option := range options {
		option(m)
	}
	args := m.Command
	if len(args) == 0 {
		args = []string{"migrate"}
	}
	env := PostgresEnvVars(database, databaseName)
	if m.AdditionalEnv != nil {
		env = append(env, m.AdditionalEnv...)
	}
	return v1.Container{
		Name:  "migrate",
		Image: common.GetImage(serviceName, version),
		Args:  args,
		Env:   env,
	}
}
