package v1beta2

import (
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/typeutils"
)

type ScalingSpec struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	MinReplica int `json:"minReplica,omitempty"`
	// +optional
	MaxReplica int `json:"maxReplica,omitempty"`
	// +optional
	CpuLimit int `json:"cpuLimit,omitempty"`
}

type DatabaseSpec struct {
	// +optional
	Url string `json:"url,omitempty"`
	// +optional
	Type string `json:"type,omitempty"`
}

type IngressConfig struct {
	// +optional
	Annotations map[string]string `json:"annotations"`
}

func (cfg *IngressConfig) Compute(stack *Stack, config *ConfigurationSpec, path string) *apisv1beta2.IngressSpec {
	annotations := make(map[string]string)
	if config.Ingress != nil && config.Ingress.Annotations != nil {
		annotations = typeutils.MergeMaps(annotations, config.Ingress.Annotations)
	}
	if cfg != nil && cfg.Annotations != nil {
		annotations = typeutils.MergeMaps(annotations, cfg.Annotations)
	}
	var ingressTLS *apisv1beta2.IngressTLS
	if config.Ingress != nil {
		ingressTLS = config.Ingress.TLS
	}
	return &apisv1beta2.IngressSpec{
		Path:        path,
		Host:        stack.Spec.Host,
		Annotations: annotations,
		TLS:         ingressTLS,
	}
}
