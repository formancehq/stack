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

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CommonServiceProperties struct{}

// ServerSpec defines the desired state of Server
type ServerSpec struct {
	CommonServiceProperties `json:",inline"`
	// +optional
	InitContainers []corev1.Container `json:"containers,omitempty"`
	// +optional
	ResourcesConfigMap string `json:"resourcesConfigMap"`
	// +optional
	TemplatesConfigMap string `json:"templatesConfigMap"`
	// +optional
	StreamsConfigMap string `json:"streamsConfigMap"`
	// +optional
	GlobalConfigMap string `json:"globalConfigMap"`
	// +optional
	Env []corev1.EnvVar `json:"env"`
	// +optional
	ConfigurationFile string `json:"configurationFile"`
}

type Status struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// Server is the Schema for the servers API
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec `json:"spec,omitempty"`
	Status Status     `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServerList contains a list of Server
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Server{}, &ServerList{})
}
