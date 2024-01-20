package services

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Create(ctx core.Context, owner v1beta1.Module, name string, mutators ...core.ObjectMutator[*corev1.Service]) (*corev1.Service, error) {
	mutators = append(mutators,
		Configure(name),
		core.WithController[*corev1.Service](ctx.GetScheme(), owner),
	)
	service, _, err := core.CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Name:      name,
		Namespace: owner.GetStack(),
	}, mutators...)
	return service, err
}

func Configure(name string) func(service *corev1.Service) {
	return func(t *corev1.Service) {
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": name,
		}
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       8080,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt32(8080),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
		}
	}
}
