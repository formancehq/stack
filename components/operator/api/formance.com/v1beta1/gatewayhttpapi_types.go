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

type GatewayHTTPAPIRule struct {
	Path string `json:"path"`
	//+optional
	Methods []string `json:"methods"`
	//+optional
	//+kubebuilder:default:=false
	Secured bool `json:"secured"`
}

type GatewayHTTPAPISpec struct {
	StackDependency `json:",inline"`
	// Name indicates prefix api
	Name string `json:"name"`
	// Rules
	Rules []GatewayHTTPAPIRule `json:"rules"`
	// Health check endpoint
	HealthCheckEndpoint string `json:"healthCheckEndpoint,omitempty"`
}

type GatewayHTTPAPIStatus struct {
	Status `json:",inline"`
	//+optional
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"

// GatewayHTTPAPI is the Schema for the HTTPAPIs API
type GatewayHTTPAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewayHTTPAPISpec   `json:"spec,omitempty"`
	Status GatewayHTTPAPIStatus `json:"status,omitempty"`
}

func (in *GatewayHTTPAPI) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *GatewayHTTPAPI) IsReady() bool {
	return in.Status.Ready
}

func (in *GatewayHTTPAPI) SetError(s string) {
	in.Status.Info = s
}

func (a GatewayHTTPAPI) GetStack() string {
	return a.Spec.Stack
}

func (in *GatewayHTTPAPI) GetConditions() *Conditions {
	return &in.Status.Conditions
}

//+kubebuilder:object:root=true

// GatewayHTTPAPIList contains a list of GatewayHTTPAPI
type GatewayHTTPAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GatewayHTTPAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GatewayHTTPAPI{}, &GatewayHTTPAPIList{})
}
