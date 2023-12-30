package deployments

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func DefaultLiveness(port string, opts ...ProbeOpts) *corev1.Probe {
	return liveness(newProbeHandler(port, opts...))
}
func liveness(handler corev1.ProbeHandler) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler:                  handler,
		InitialDelaySeconds:           1,
		TimeoutSeconds:                30,
		PeriodSeconds:                 2,
		SuccessThreshold:              1,
		FailureThreshold:              20,
		TerminationGracePeriodSeconds: pointer.Int64(10),
	}
}

type ProbeOpts func(*corev1.ProbeHandler) *corev1.ProbeHandler

func newProbeHandler(port string, opts ...ProbeOpts) corev1.ProbeHandler {
	probe := corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Port:   intstr.FromString(port),
			Scheme: "HTTP",
		},
	}

	for _, opt := range append(defaultProbeOptions, opts...) {
		opt(&probe)
	}

	return probe
}

var defaultProbeOptions = []ProbeOpts{
	WithProbePath("/_healthcheck"),
}

func WithHost(host string) ProbeOpts {
	return func(p *corev1.ProbeHandler) *corev1.ProbeHandler {
		p.HTTPGet.Host = host
		return p
	}
}

func WithProbePath(path string) ProbeOpts {
	return func(p *corev1.ProbeHandler) *corev1.ProbeHandler {
		p.HTTPGet.Path = path
		return p
	}
}
