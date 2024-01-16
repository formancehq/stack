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

// VersionsHistorySpec defines the desired state of VersionsHistory
type VersionsHistorySpec struct {
	StackDependency `json:",inline"`
	Module          string `json:"module"`
	Version         string `json:"version"`
}

// VersionsHistoryStatus defines the observed state of VersionsHistory
type VersionsHistoryStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"

// VersionsHistory is the Schema for the versionshistories API
type VersionsHistory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VersionsHistorySpec   `json:"spec,omitempty"`
	Status VersionsHistoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VersionsHistoryList contains a list of VersionsHistory
type VersionsHistoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VersionsHistory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VersionsHistory{}, &VersionsHistoryList{})
}
