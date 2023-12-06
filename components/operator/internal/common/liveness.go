package common

import (
	"net/url"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func DefaultLiveness(port int32) *corev1.Probe {
	return liveness(newProbeHandler(port, "/_healthcheck"))
}

func LegacyLiveness(port int32) *corev1.Probe {
	return liveness(newProbeHandler(port, "/_health"))
}

func LivenessEndpoint(str string) *corev1.Probe {

	//str as url path
	url, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	port, err := strconv.ParseInt(url.Port(), 10, 64)
	if err != nil {
		panic(err)
	}

	return liveness(
		newProbeHandler(
			int32(port),
			url.Path,
			withHost(url.Hostname()),
		),
	)
}

type ProbeOpts func(*corev1.ProbeHandler) *corev1.ProbeHandler

func newProbeHandler(port int32, path string, opts ...ProbeOpts) corev1.ProbeHandler {
	probe := corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: path,
			Port: intstr.IntOrString{
				IntVal: port,
			},
			Scheme: "HTTP",
		},
	}

	for _, opt := range opts {
		opt(&probe)
	}

	return probe
}

func withHost(host string) ProbeOpts {
	return func(p *corev1.ProbeHandler) *corev1.ProbeHandler {
		p.HTTPGet.Host = host
		return p
	}
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
