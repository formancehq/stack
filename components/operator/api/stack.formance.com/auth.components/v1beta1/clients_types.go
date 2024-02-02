/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientConfiguration struct {
	// +optional
	Public bool `json:"public"`
	// +optional
	Description *string `json:"description,omitempty"`
	// +optional
	RedirectUris []string `json:"redirectUris,omitempty"`
	// +optional
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty"`
	// +optional
	Scopes []string `json:"scopes,omitempty"`
}

type StaticClient struct {
	ClientConfiguration `json:",inline" yaml:",inline"`
	ID                  string `json:"id" yaml:"id"`
	// +optional
	Secrets []string `json:"secrets" yaml:"secrets"`
}

// ClientSpec defines the desired state of Client
type ClientSpec struct {
	ClientConfiguration `json:",inline"`
	AuthServerReference string `json:"authServerReference"`
}

// ClientStatus defines the observed state of Client
type ClientStatus struct {
	Status       `json:",inline"`
	AuthServerID string `json:"authServerID,omitempty"`
	// +optional
	Scopes map[string]string `json:"scopes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Server ID",type="string",JSONPath=".status.authServerID",description="Auth server ID"

// Client is the Schema for the oauths API
type Client struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClientSpec   `json:"spec,omitempty"`
	Status ClientStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClientList contains a list of Client
type ClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Client `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Client{}, &ClientList{})
}
