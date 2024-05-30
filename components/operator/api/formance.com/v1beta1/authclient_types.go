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

type AuthClientSpec struct {
	StackDependency `json:",inline" yaml:",inline"`
	// ID indicates the client id
	// It must be used with oauth2 `client_id` parameter
	ID string `json:"id" yaml:"id"`
	// +optional
	// Public indicate whether a client is confidential or not.
	// Confidential clients are clients which the secret can be kept secret...
	// As opposed to public clients which cannot have a secret (application single page for example)
	// +kubebuilder:default:=false
	Public bool `json:"public" yaml:"public"`
	// +optional
	// Description represents an optional description of the client
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// +optional
	// RedirectUris allow to list allowed redirect uris for the client
	RedirectUris []string `json:"redirectUris,omitempty" yaml:"redirectUris"`
	// +optional
	// RedirectUris allow to list allowed post logout redirect uris for the client
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty" yaml:"PostLogoutRedirectUris"`
	// +optional
	// Scopes allow to five some scope to the client
	Scopes []string `json:"scopes,omitempty" yaml:"scopes"`
	// +optional
	// Secret allow to configure a secret for the client.
	// It is not required as some client could use some oauth2 flows which does not requires a client secret
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

type AuthClientStatus struct {
	Status `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"

// AuthClient allow to create OAuth2/OIDC clients on the auth server (see [Auth](#auth))
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

func (in *AuthClient) GetConditions() *Conditions {
	return &in.Status.Conditions
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
