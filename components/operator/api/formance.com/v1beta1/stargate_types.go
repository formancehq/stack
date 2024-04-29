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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StargateAuthSpec struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Issuer       string `json:"issuer"`
}

// StargateSpec defines the desired state of Stargate
type StargateSpec struct {
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
	ServerURL        string           `json:"serverURL"`
	OrganizationID   string           `json:"organizationID"`
	StackID          string           `json:"stackID"`
	Auth             StargateAuthSpec `json:"auth"`
}

// StargateStatus defines the observed state of Stargate
type StargateStatus struct {
	ModuleStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
// +kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
// +kubebuilder:metadata:labels=formance.com/kind=module
// Stargate is the Schema for the stargates API
type Stargate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StargateSpec   `json:"spec,omitempty"`
	Status StargateStatus `json:"status,omitempty"`
}

func (in *Stargate) IsEE() bool {
	return false
}

func (in *Stargate) GetVersion() string {
	return in.Spec.Version
}

func (in *Stargate) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Stargate) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Stargate) IsReady() bool {
	return in.Status.Ready
}

func (in *Stargate) SetError(s string) {
	in.Status.Info = s
}

func (a Stargate) IsDebug() bool {
	return a.Spec.Debug
}

func (a Stargate) IsDev() bool {
	return a.Spec.Dev
}

func (s Stargate) GetStack() string {
	return s.Spec.Stack
}

func (s *Stargate) SetCondition(condition Condition) {
	s.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// StargateList contains a list of Stargate
type StargateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stargate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stargate{}, &StargateList{})
}
