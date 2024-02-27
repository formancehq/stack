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

// OrchestrationSpec defines the desired state of Orchestration
type OrchestrationSpec struct {
	StackDependency  `json:",inline"`
	ModuleProperties `json:",inline"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// OrchestrationStatus defines the observed state of Orchestration
type OrchestrationStatus struct {
	ModuleStatus `json:",inline"`
	//+optional
	TemporalURI *URI `json:"temporalURI,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Orchestration is the Schema for the orchestrations API
type Orchestration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrchestrationSpec   `json:"spec,omitempty"`
	Status OrchestrationStatus `json:"status,omitempty"`
}

func (in *Orchestration) isEventPublisher() {}

func (in *Orchestration) IsEE() bool {
	return false
}

func (in *Orchestration) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Orchestration) IsReady() bool {
	return in.Status.Ready
}

func (in *Orchestration) SetError(s string) {
	in.Status.Info = s
}

func (in *Orchestration) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Orchestration) GetVersion() string {
	return in.Spec.Version
}

func (a Orchestration) GetStack() string {
	return a.Spec.Stack
}

func (a *Orchestration) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

func (a Orchestration) IsDebug() bool {
	return a.Spec.Debug
}

func (a Orchestration) IsDev() bool {
	return a.Spec.Dev
}

//+kubebuilder:object:root=true

// OrchestrationList contains a list of Orchestration
type OrchestrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Orchestration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Orchestration{}, &OrchestrationList{})
}

var _ EventPublisher = (*Orchestration)(nil)
