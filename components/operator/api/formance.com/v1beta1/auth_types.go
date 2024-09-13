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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DelegatedOIDCServerConfiguration struct {
	// Issuer is the url of the delegated oidc server
	Issuer string `json:"issuer,omitempty"`
	// ClientID is the client id to use for authentication
	ClientID string `json:"clientID,omitempty"`
	// ClientSecret is the client secret to use for authentication
	ClientSecret string `json:"clientSecret,omitempty"`
}

type AuthSpec struct {
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
	//+optional
	// Contains information about a delegated authentication server to use to delegate authentication
	DelegatedOIDCServer *DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer,omitempty"`
	//+optional
	// Allow to override the default signing key used to sign JWT tokens.
	SigningKey string `json:"signingKey,omitempty"`
	//+optional
	// Allow to override the default signing key used to sign JWT tokens using a k8s secret
	SigningKeyFromSecret *v1.SecretKeySelector `json:"signingKeyFromSecret,omitempty"`
	//+optional
	// Allow to enable scopes usage on authentication.
	//
	// If not enabled, each service will check the authentication but will not restrict access following scopes.
	// in this case, if authenticated, it is ok.
	// +kubebuilder:default:=false
	EnableScopes bool `json:"enableScopes"`
}

type AuthStatus struct {
	Status `json:",inline"`
	//+optional
	// Clients contains the list of clients created using [AuthClient](#authclient)
	Clients []string `json:"clients"`
}

// Auth represent the authentication module of a stack.
//
// It is an OIDC compliant server.
//
// Creating it for a stack automatically add authentication on all supported modules.
//
// The auth service is basically a proxy to another OIDC compliant server.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Clients",type=string,JSONPath=".status.clients",description="Synchronized auth clients"
// +kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
// +kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
// +kubebuilder:metadata:labels=formance.com/kind=module
type Auth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthSpec   `json:"spec,omitempty"`
	Status AuthStatus `json:"status,omitempty"`
}

func (in *Auth) IsEE() bool {
	return false
}

func (in *Auth) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Auth) IsReady() bool {
	return in.Status.Ready
}

func (in *Auth) SetError(s string) {
	in.Status.Info = s
}

func (in *Auth) GetConditions() *Conditions {
	return &in.Status.Conditions
}

func (in *Auth) GetVersion() string {
	return in.Spec.Version
}

func (in Auth) GetStack() string {
	return in.Spec.Stack
}

func (in Auth) IsDebug() bool {
	return in.Spec.Debug
}

func (in Auth) IsDev() bool {
	return in.Spec.Dev
}

//+kubebuilder:object:root=true

// AuthList contains a list of Auth
type AuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Auth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Auth{}, &AuthList{})
}
