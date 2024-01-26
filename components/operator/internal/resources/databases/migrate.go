package databases

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	v1 "k8s.io/api/core/v1"
)

type MigrationConfiguration struct {
	Command       []string
	AdditionalEnv []v1.EnvVar
}

func MigrateDatabaseContainer(image string, database *v1beta1.Database, options ...func(m *MigrationConfiguration)) v1.Container {
	m := &MigrationConfiguration{}
	for _, option := range options {
		option(m)
	}
	args := m.Command
	if len(args) == 0 {
		args = []string{"migrate"}
	}

	env := GetPostgresEnvVars(database)
	if m.AdditionalEnv != nil {
		env = append(env, m.AdditionalEnv...)
	}

	return v1.Container{
		Name:  "migrate",
		Image: image,
		Args:  args,
		Env:   env,
	}
}
