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

type TemporalTLSConfig struct {
	// +optional
	CRT string `json:"crt"`
	// +optional
	Key string `json:"key"`
	// +optional
	SecretName string `json:"secretName"`
}

type TemporalConfig struct {
	Address   string            `json:"address"`
	Namespace string            `json:"namespace"`
	TLS       TemporalTLSConfig `json:"tls,omitempty"`
}

// OrchestrationSpec defines the desired state of Orchestration
type OrchestrationSpec struct {
	StackDependency         `json:",inline"`
	CommonServiceProperties `json:",inline"`
	Temporal                TemporalConfig `json:"temporal"`
}

// OrchestrationStatus defines the observed state of Orchestration
type OrchestrationStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// Orchestration is the Schema for the orchestrations API
type Orchestration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrchestrationSpec   `json:"spec,omitempty"`
	Status OrchestrationStatus `json:"status,omitempty"`
}

func (a Orchestration) GetStack() string {
	return a.Spec.Stack
}

func (a *Orchestration) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
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
