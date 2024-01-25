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
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AuthClientSpec defines the desired state of AuthClient
type AuthClientSpec struct {
	StackDependency `json:",inline" yaml:",inline"`
	ID              string `json:"id" yaml:"id"`
	// +optional
	Public bool `json:"public" yaml:"public"`
	// +optional
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// +optional
	RedirectUris []string `json:"redirectUris,omitempty" yaml:"redirectUris"`
	// +optional
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty" yaml:"PostLogoutRedirectUris"`
	// +optional
	Scopes []string `json:"scopes,omitempty" yaml:"scopes"`
	// +optional
	Secret string `json:"secret,omitempty"`
}

var _ yaml.Marshaler = (*AuthClientSpec)(nil)

func (spec AuthClientSpec) MarshalYAML() (interface{}, error) {
	type aux AuthClientSpec
	return struct {
		aux     `yaml:",inline"`
		Secrets []string `yaml:"secrets"`
	}{
		aux: aux(spec),
		Secrets: func() []string {
			if spec.Secret == "" {
				return []string{}
			}
			return []string{spec.Secret}
		}(),
	}, nil
}

// AuthClientStatus defines the observed state of AuthClient
type AuthClientStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"

// AuthClient is the Schema for the authclients API
type AuthClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthClientSpec   `json:"spec,omitempty"`
	Status AuthClientStatus `json:"status,omitempty"`
}

func (a *AuthClient) SetReady(b bool) {
	a.Status.Ready = b
}

func (in *AuthClient) IsReady() bool {
	return in.Status.Ready
}

func (a *AuthClient) SetError(s string) {
	a.Status.Info = s
}

func (a AuthClient) GetStack() string {
	return a.Spec.Stack
}

//+kubebuilder:object:root=true

// AuthClientList contains a list of AuthClient
type AuthClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AuthClient{}, &AuthClientList{})
}
