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

// DatabaseConfigurationSpec defines the desired state of DatabaseConfiguration
type DatabaseConfigurationSpec struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	// +optional
	Username string `json:"username"`
	// +optional
	Password string `json:"password"`

	// +optional
	Debug bool `json:"debug"`
	// +optional
	CredentialsFromSecret string `json:"credentialsFromSecret"`
	// +optional
	DisableSSLMode bool `json:"disableSSLMode"`
}

// DatabaseConfigurationStatus defines the observed state of DatabaseConfiguration
type DatabaseConfigurationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// DatabaseConfiguration is the Schema for the databaseconfigurations API
type DatabaseConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseConfigurationSpec   `json:"spec,omitempty"`
	Status DatabaseConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseConfigurationList contains a list of DatabaseConfiguration
type DatabaseConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatabaseConfiguration{}, &DatabaseConfigurationList{})
}
