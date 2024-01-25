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

// ElasticSearchConfiguration defines the desired state of ElasticSearchConfiguration
type ElasticSearchConfiguration struct {
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

func (in *ElasticSearchConfiguration) Endpoint() string {
	return fmt.Sprintf("%s://%s:%d", in.Scheme, in.Host, in.Port)
}

// SearchSpec defines the desired state of Search
type SearchSpec struct {
	StackDependency  `json:",inline"`
	ModuleProperties `json:",inline"`
	//+optional
	Batching *Batching `json:"batching,omitempty"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// SearchStatus defines the observed state of Search
type SearchStatus struct {
	ModuleStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Search is the Schema for the searches API
type Search struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchSpec   `json:"spec,omitempty"`
	Status SearchStatus `json:"status,omitempty"`
}

func (in *Search) IsEE() bool {
	return false
}

func (in *Search) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Search) IsReady() bool {
	return in.Status.Ready
}

func (in *Search) SetError(s string) {
	in.Status.Info = s
}

func (in *Search) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Search) GetVersion() string {
	return in.Spec.Version
}

func (a Search) GetStack() string {
	return a.Spec.Stack
}

func (a *Search) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

func (a Search) IsDebug() bool {
	return a.Spec.Debug
}

func (a Search) IsDev() bool {
	return a.Spec.Dev
}

//+kubebuilder:object:root=true

// SearchList contains a list of Search
type SearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Search `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Search{}, &SearchList{})
}
