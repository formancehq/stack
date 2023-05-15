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

package v1beta2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type StargateSpec struct {
	CommonServiceProperties `json:",inline"`
	Scalable                `json:",inline"`

	// +optional
	Monitoring *MonitoringSpec `json:"monitoring"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// Stargate is the Schema for the stargate client service
type Stargate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StargateSpec `json:"spec,omitempty"`
	Status Status       `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StargateList contains a list of Stargate Clients
type StargateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stargate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stargate{}, &StargateList{})
}
