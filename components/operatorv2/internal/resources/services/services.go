package services

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ConfigureK8SService(name string, options ...func(service *v1.Service)) func(service *v1.Service) {
	return func(t *v1.Service) {
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": name,
		}
		t.Spec = v1.ServiceSpec{
			Ports: []v1.ServicePort{{
				Name:       "http",
				Port:       8080,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt32(8080),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
		}
		for _, option := range options {
			option(t)
		}
	}
}
