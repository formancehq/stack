package v1beta3

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
)

type ElasticSearchTLSConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	SkipCertVerify bool `json:"skipCertVerify,omitempty"`
}

type ElasticSearchBasicAuthConfig struct {
	// +optional
	Username string `json:"username"`
	// +optional
	Password string `json:"password"`
	// +optional
	SecretName string `json:"secretName"`
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

type ElasticSearchConfigurationSpec struct {
	v1beta1.ElasticSearchConfigurationSpec `json:",inline"`
	PathPrefix                             string `json:"pathPrefix"`
}

// +kubebuilder:object:generate=true
type SearchSpec struct {
	CommonServiceProperties `json:",inline"`
	ElasticSearchConfig     ElasticSearchConfigurationSpec `json:"elasticSearch"`

	// +optional
	Batching v1beta1.Batching `json:"batching"`

	// +optional
	SearchResourceProperties *ResourceProperties `json:"searchResourceProperties,omitempty"`
	// +optional
	BenthosResourceProperties *ResourceProperties     `json:"benthosResourceProperties,omitempty"`
	Annotations               AnnotationsServicesSpec `json:"annotations,omitempty"`
}
