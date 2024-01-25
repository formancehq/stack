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
	Issuer       string `json:"issuer,omitempty"`
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

// AuthSpec defines the desired state of Auth
type AuthSpec struct {
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
	//+optional
	DelegatedOIDCServer *DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer,omitempty"`
	//+optional
	SigningKey string `json:"signingKey,omitempty"`
	//+optional
	SigningKeyFromSecret *v1.SecretKeySelector `json:"signingKeyFromSecret,omitempty"`
	//+optional
	EnableScopes bool `json:"enableScopes"`
}

// AuthStatus defines the observed state of Auth
type AuthStatus struct {
	ModuleStatus `json:",inline"`
	//+optional
	Clients []string `json:"clients"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Clients",type=string,JSONPath=".status.clients",description="Synchronized auth clients"
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Auth is the Schema for the auths API
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

func (in *Auth) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Auth) GetVersion() string {
	return in.Spec.Version
}

func (a Auth) GetStack() string {
	return a.Spec.Stack
}

func (a *Auth) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

func (a Auth) IsDebug() bool {
	return a.Spec.Debug
}

func (a Auth) IsDev() bool {
	return a.Spec.Dev
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
