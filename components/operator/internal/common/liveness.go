package common

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func DefaultLiveness(port int32) *corev1.Probe {
	return liveness(port, "/_healthcheck")
}

func LegacyLiveness(port int32) *corev1.Probe {
	return liveness(port, "/_health")
}

func liveness(port int32, path string) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.IntOrString{
					IntVal: port,
				},
				Scheme: "HTTP",
			},
		},
		InitialDelaySeconds:           1,
		TimeoutSeconds:                30,
		PeriodSeconds:                 2,
		SuccessThreshold:              1,
		FailureThreshold:              10,
		TerminationGracePeriodSeconds: pointer.Int64(10),
	}
}
