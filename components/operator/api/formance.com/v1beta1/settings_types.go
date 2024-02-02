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

// SettingsSpec defines the desired state of Settings
type SettingsSpec struct {
	//+optional
	Stacks []string `json:"stacks,omitempty"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
}

// SettingsStatus defines the observed state of Settings
type SettingsStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Key",type=string,JSONPath=".spec.key",description="Key"
//+kubebuilder:printcolumn:name="Value",type=string,JSONPath=".spec.value",description="Value"

// Settings is the Schema for the settings API
type Settings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SettingsSpec   `json:"spec,omitempty"`
	Status SettingsStatus `json:"status,omitempty"`
}

func (in *Settings) GetStacks() []string {
	return in.Spec.Stacks
}

func (in *Settings) IsWildcard() bool {
	return len(in.Spec.Stacks) == 1 && in.Spec.Stacks[0] == "*"
}

//+kubebuilder:object:root=true

// SettingsList contains a list of Settings
type SettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Settings `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Settings{}, &SettingsList{})
}
