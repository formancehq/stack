package services

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Create(ctx core.Context, owner v1beta1.Dependent, name string, mutators ...core.ObjectMutator[*corev1.Service]) (*corev1.Service, error) {
	mutators = append(mutators,
		core.WithController[*corev1.Service](ctx.GetScheme(), owner),
	)
	mutators = append(mutators, func(t *corev1.Service) error {
		var err error
		t.ObjectMeta.Annotations, err = settings.GetMapOrEmpty(ctx, owner.GetStack(), "services", name, "annotations")
		if err != nil {
			return err
		}
		return nil
	})
	service, _, err := core.CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Name:      name,
		Namespace: owner.GetStack(),
	}, mutators...)
	return service, err
}

func WithDefault(name string) core.ObjectMutator[*corev1.Service] {
	return func(t *corev1.Service) error {
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": name,
		}
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       8080,
				Protocol:   "TCP",
				TargetPort: intstr.FromString("http"),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
		}

		return nil
	}
}
