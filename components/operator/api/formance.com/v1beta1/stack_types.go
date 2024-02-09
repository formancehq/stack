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

// StackSpec defines the desired state of Stack
// The version of the stack can be specified using either the field `version` or the `versionsFromFile` field.
// The `version` field will have priority over `versionFromFile`
// If `versions` and `versionsFromFile` are not specified, "latest" will be used.
type StackSpec struct {
	DevProperties `json:",inline"`
	// +optional
	// Version allow to specify the version of the components
	// Must be a valid docker tag
	Version string `json:"version,omitempty"`
	// +optional
	// VersionsFromFile allow to specify a formance.com/Versions object which contains individual versions
	// for each component.
	// Must reference a valid formance.com/Versions object
	VersionsFromFile string `json:"versionsFromFile"`
	//+optional
	EnableAudit bool `json:"enableAudit,omitempty"`
	//+optional
	Disabled bool `json:"disabled"`
}

// StackStatus defines the observed state of Stack
type StackStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Disable",type=string,JSONPath=".spec.disabled",description="Stack Disabled"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Stack Version"
//+kubebuilder:printcolumn:name="Versions From file",type="string",JSONPath=".spec.versionsFromFile",description="Stack Version From File"
//+kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready",description="Is stack ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Stack is the Schema for the stacks API
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

func (in *Stack) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *Stack) IsReady() bool {
	return in.Status.Ready
}

func (in *Stack) SetError(s string) {
	in.Status.SetError(s)
}

func (in *Stack) GetVersion() string {
	if in.Spec.Version == "" {
		return "latest"
	}
	return in.Spec.Version
}

//+kubebuilder:object:root=true

// StackList contains a list of Stack
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stack{}, &StackList{})
}
