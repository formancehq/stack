package v1beta3

import (
	"fmt"
)

type Batching struct {
	Count  int    `json:"count"`
	Period string `json:"period"`
}

type ElasticSearchTLSConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	SkipCertVerify bool `json:"skipCertVerify,omitempty"`
}

type ElasticSearchBasicAuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ElasticSearchConfig struct {
	// +optional
	// +kubebuilder:validation:Enum:={http,https}
	// +kubebuilder:validation:default:=https
	Scheme string `json:"scheme,omitempty"`
	Host   string `json:"host,omitempty"`
	Port   uint16 `json:"port,omitempty"`
	// +optional
	TLS ElasticSearchTLSConfig `json:"tls"`
	// +optional
	BasicAuth *ElasticSearchBasicAuthConfig `json:"basicAuth,omitempty"`
	// +optional
	PathPrefix string `json:"pathPrefix"`
	// +optional
	UseZinc bool `json:"useZinc,omitempty"`
}

func (in *ElasticSearchConfig) Endpoint() string {
	return fmt.Sprintf("%s://%s:%d%s", in.Scheme, in.Host, in.Port, in.PathPrefix)
}

// +kubebuilder:object:generate=true
type SearchSpec struct {
	ElasticSearchConfig ElasticSearchConfig `json:"elasticSearch"`

	// +optional
	Batching Batching `json:"batching"`

	// +optional
	Annotations AnnotationsServicesSpec `json:"service"`
}

const DefaultESIndex = "stacks"
