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
	// +optional
	// +kubebuilder:default:=false
	// EnableAudit enable audit at the stack level.
	// Actually, it enables audit on [Gateway](#gateway)
	EnableAudit bool `json:"enableAudit,omitempty"`
	// +optional
	// +kubebuilder:default:=false
	// Disabled indicate the stack is disabled.
	// A disabled stack disable everything
	// It just keeps the namespace and the [Database](#database) resources.
	Disabled bool `json:"disabled"`
}

type StackStatus struct {
	Status `json:",inline"`
	// Modules register detected modules
	Modules []string `json:"modules,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Disable",type=string,JSONPath=".spec.disabled",description="Stack Disabled"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Stack Version"
//+kubebuilder:printcolumn:name="Versions From file",type="string",JSONPath=".spec.versionsFromFile",description="Stack Version From File"
//+kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready",description="Is stack ready"
//+kubebuilder:printcolumn:name="Modules",type=string,JSONPath=".status.modules",description="Modules List Registered"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Stack represents a formance stack.
// A Stack is basically a container. It holds some global properties and
// creates a namespace if not already existing.
//
// To do more, you need to create some [modules](#modules).
//
// The Stack resource allow to specify the version of the stack.
//
// It can be specified using either the field `.spec.version` or the `.spec.versionsFromFile` field (Refer to the documentation of [Versions](#versions) resource.
//
// The `version` field will have priority over `versionFromFile`.
//
// If `versions` and `versionsFromFile` are not specified, "latest" will be used.
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

func (in *Stack) MustSkip() bool {
	return in.GetAnnotations()[SkipLabel] == "true"
}

func (in *Stack) GetConditions() *Conditions {
	return &in.Status.Conditions
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
