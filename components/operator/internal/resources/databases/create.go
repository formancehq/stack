package databases

import (
	"strings"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx core.Context, owner interface {
	client.Object
	SetCondition(condition v1beta1.Condition)
	GetStack() string
}) (*v1beta1.Database, error) {
	condition := v1beta1.Condition{
		Type:               "DatabaseReady",
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
		func(t *v1beta1.Database) error {
			t.Spec.Stack = owner.GetStack()
			t.Spec.Service = serviceName

			return nil
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
	} else {
		condition.Message = "database is ok"
		condition.Status = metav1.ConditionTrue
	}

	return database, err
}
