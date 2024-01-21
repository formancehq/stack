/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// ElasticSearchConfigurationSpec defines the desired state of ElasticSearchConfiguration
type ElasticSearchConfigurationSpec struct {
	ConfigurationProperties `json:",inline"`
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
}

func (in *ElasticSearchConfigurationSpec) Endpoint() string {
	return fmt.Sprintf("%s://%s:%d", in.Scheme, in.Host, in.Port)
}

// ElasticSearchConfigurationStatus defines the observed state of ElasticSearchConfiguration
type ElasticSearchConfigurationStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ElasticSearchConfiguration is the Schema for the elasticsearchconfigs API
type ElasticSearchConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticSearchConfigurationSpec   `json:"spec,omitempty"`
	Status ElasticSearchConfigurationStatus `json:"status,omitempty"`
}

func (in *ElasticSearchConfiguration) GetStacks() []string {
	return in.Spec.Stacks
}

func (in *ElasticSearchConfiguration) IsWildcard() bool {
	return in.Spec.ApplyOnAllStacks
}

var _ ConfigurationObject = &ElasticSearchConfiguration{}

//+kubebuilder:object:root=true

// ElasticSearchConfigurationList contains a list of ElasticSearchConfiguration
type ElasticSearchConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticSearchConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticSearchConfiguration{}, &ElasticSearchConfigurationList{})
}
