package databases

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	v1 "k8s.io/api/core/v1"
)

type MigrationConfiguration struct {
	Command       []string
	AdditionalEnv []v1.EnvVar
}

func MigrateDatabaseContainer(ctx core.Context, stack *v1beta1.Stack, image string, database *v1beta1.Database, options ...func(m *MigrationConfiguration)) (v1.Container, error) {
	m := &MigrationConfiguration{}
	for _, option := range options {
		option(m)
	}
	args := m.Command
	if len(args) == 0 {
		args = []string{"migrate"}
	}

	env, err := GetPostgresEnvVars(ctx, stack, database)
	if err != nil {
		return v1.Container{}, err
	}

	if m.AdditionalEnv != nil {
		env = append(env, m.AdditionalEnv...)
	}

	return v1.Container{
		Name:  "migrate",
		Image: image,
		Args:  args,
		Env:   env,
	}, nil
}

func Migrate(ctx core.Context, stack *v1beta1.Stack, image string, database *v1beta1.Database, options ...func(m *MigrationConfiguration)) error {
	container, err := MigrateDatabaseContainer(ctx, stack, image, database, options...)
	if err != nil {
		return err
	}

	return jobs.Handle(ctx, database, fmt.Sprintf("%s-migration", database.Name), container)
}
