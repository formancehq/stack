package versions

import (
	"reflect"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
)

// +kubebuilder:rbac:groups=stack.formance.com,resources=versions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions/finalizers,verbs=update

func Reconcile(ctx Context, req *v1beta3.Versions) error {
	_, _, err := CreateOrUpdate[*v1beta1.Versions](ctx, types.NamespacedName{
		Name: req.Name,
	}, func(t *v1beta1.Versions) error {
		t.Spec = map[string]string{}

		for i := 0; i < reflect.ValueOf(req.Spec).NumField(); i++ {
			field := reflect.TypeOf(req.Spec).Field(i)
			jsonTag := field.Tag.Get("json")
			name := strings.Split(jsonTag, ",")[0]
			value := reflect.ValueOf(req.Spec).Field(i).Interface().(string)
			if value == "" {
				continue
			}
			t.Spec[name] = value
		}

		return nil
	})
	return err
}

func init() {
	Init(
		WithReconciler[*v1beta3.Versions](Reconcile),
	)
}
