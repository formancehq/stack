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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HTTPAPIRule struct {
	Path string `json:"path"`
	//+optional
	Methods []string `json:"methods"`
	//+optional
	Secured bool `json:"secured"`
}

// HTTPAPISpec defines the desired state of HTTPAPI
type HTTPAPISpec struct {
	StackDependency `json:",inline"`
	// Name indicates prefix api
	Name string `json:"name"`
	//+optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Rules
	Rules []HTTPAPIRule `json:"rules"`
}

// HTTPAPIStatus defines the observed state of HTTPAPI
type HTTPAPIStatus struct {
	CommonStatus `json:",inline"`
	//+optional
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"

// HTTPAPI is the Schema for the HTTPAPIs API
type HTTPAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HTTPAPISpec   `json:"spec,omitempty"`
	Status HTTPAPIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HTTPAPIList contains a list of HTTPAPI
type HTTPAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HTTPAPI{}, &HTTPAPIList{})
}
