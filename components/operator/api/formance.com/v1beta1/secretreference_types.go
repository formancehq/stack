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

// SecretReferenceSpec defines the desired state of SecretReference
type SecretReferenceSpec struct {
	StackDependency `json:",inline"`
	SecretName      string `json:"secretName"`
}

// SecretReferenceStatus defines the observed state of SecretReference
type SecretReferenceStatus struct {
	CommonStatus `json:",inline"`
	// Hash of the secret, allow to watch secrets and reload if needed
	//+optional
	Hash string `json:"hash,omitempty"`
	// Synced secret at last reconciliation
	//+optional
	SyncedSecret string `json:"syncedSecret,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// SecretReference is the Schema for the secretreferences API
type SecretReference struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretReferenceSpec   `json:"spec,omitempty"`
	Status SecretReferenceStatus `json:"status,omitempty"`
}

func (in *SecretReference) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *SecretReference) IsReady() bool {
	return in.Status.Ready
}

func (in *SecretReference) SetError(s string) {
	in.Status.SetError(s)
}

func (in *SecretReference) GetStack() string {
	return in.Spec.Stack
}

//+kubebuilder:object:root=true

// SecretReferenceList contains a list of SecretReference
type SecretReferenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretReference `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretReference{}, &SecretReferenceList{})
}
