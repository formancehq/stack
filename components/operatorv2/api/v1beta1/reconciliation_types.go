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

// ReconciliationSpec defines the desired state of Reconciliation
type ReconciliationSpec struct {
	StackDependency         `json:",inline"`
	CommonServiceProperties `json:",inline"`
}

// ReconciliationStatus defines the observed state of Reconciliation
type ReconciliationStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// Reconciliation is the Schema for the reconciliations API
type Reconciliation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReconciliationSpec   `json:"spec,omitempty"`
	Status ReconciliationStatus `json:"status,omitempty"`
}

func (a Reconciliation) GetStack() string {
	return a.Spec.Stack
}

func (a *Reconciliation) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// ReconciliationList contains a list of Reconciliation
type ReconciliationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Reconciliation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Reconciliation{}, &ReconciliationList{})
}
