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

type Path struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
	Name    string   `json:"name"`
	Secured bool     `json:"secured"`
}

type Liveness int

const (
	LivenessDefault = iota
	LivenessLegacy
	LivenessDisable
)

func (l Liveness) String() string {
	switch l {
	case LivenessLegacy:
		return "_health"
	case LivenessDefault:
		return "_healthcheck"
	}
	return ""
}

// HTTPAPISpec defines the desired state of HTTPAPI
type HTTPAPISpec struct {
	Stack string `json:"stack"`
	// Name indicates prefix api
	Name string `json:"name"`
	// Secured indicate if the service is able to handle security
	Secured bool `json:"secured"`
	// HasVersionEndpoint indicates if the service has a /_info endpoint
	HasVersionEndpoint bool `json:"hasVersionEndpoint"`
	// Liveness indicates if the service has a /_health(check) endpoint
	Liveness Liveness `json:"liveness"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Port name of the container
	PortName string `json:"portName"`
}

// HTTPAPIStatus defines the observed state of HTTPAPI
type HTTPAPIStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

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
