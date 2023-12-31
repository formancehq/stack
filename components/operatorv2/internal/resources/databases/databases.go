package databases

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func Create(ctx core.Context, owner interface {
	client.Object
	SetCondition(condition v1beta1.Condition)
	GetStack() string
}) (*v1beta1.Database, error) {
	condition := v1beta1.Condition{
		Type:               "DatabaseCreated",
		ObservedGeneration: owner.GetGeneration(),
		LastTransitionTime: metav1.Now(),
	}
	defer func() {
		owner.SetCondition(condition)
	}()

	serviceName := strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind)
	database, _, err := core.CreateOrUpdate[*v1beta1.Database](ctx, types.NamespacedName{
		Name: core.GetObjectName(owner.GetStack(), serviceName),
	},
		func(t *v1beta1.Database) {
			t.Spec.Stack = owner.GetStack()
			t.Spec.Service = serviceName
		},
		core.WithController[*v1beta1.Database](ctx.GetScheme(), owner),
	)
	if err != nil {
		condition.Message = err.Error()
		condition.Status = metav1.ConditionFalse
		return nil, err
	}
	if !database.Status.Ready {
		condition.Message = "database creation pending"
		condition.Status = metav1.ConditionFalse
		return nil, core.ErrPending
	}
	condition.Message = "database is ok"
	condition.Status = metav1.ConditionTrue

	return database, err
}

type MigrationConfiguration struct {
	Command       []string
	AdditionalEnv []v1.EnvVar
}

func MigrateDatabaseContainer(image string, database v1beta1.DatabaseConfigurationSpec, databaseName string, options ...func(m *MigrationConfiguration)) v1.Container {
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
		Image: image,
		Args:  args,
		Env:   env,
	}
}
