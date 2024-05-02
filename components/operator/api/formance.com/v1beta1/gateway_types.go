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

type GatewayIngressTLS struct {
	SecretName string `json:"secretName"`
}

type GatewayIngress struct {
	//+required
	Host string `json:"host"`
	// +kubebuilder:default:="https"
	Scheme      string            `json:"scheme"`
	Annotations map[string]string `json:"annotations,omitempty"`
	//+optional
	TLS *GatewayIngressTLS `json:"tls,omitempty"`
}

// GatewaySpec defines the desired state of Gateway
type GatewaySpec struct {
	StackDependency  `json:",inline"`
	ModuleProperties `json:",inline"`
	//+optional
	Ingress *GatewayIngress `json:"ingress,omitempty"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	ModuleStatus `json:",inline"`
	//+optional
	SyncHTTPAPIs []string `json:"syncHTTPAPIs"`
	// +kubebuilder:default:=false
	AuthEnabled bool `json:"authEnabled"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="HTTP APIs",type=string,JSONPath=".status.syncHTTPAPIs",description="Synchronized http apis"
//+kubebuilder:printcolumn:name="Auth enabled",type=string,JSONPath=".status.authEnabled",description="Is authentication enabled"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
//+kubebuilder:metadata:labels=formance.com/kind=module

// Gateway is the Schema for the gateways API
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

func (in *Gateway) IsEE() bool {
	return false
}

func (in *Gateway) GetVersion() string {
	return in.Spec.Version
}

func (in *Gateway) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Gateway) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Gateway) IsReady() bool {
	return in.Status.Ready
}

func (in *Gateway) SetError(s string) {
	in.Status.Info = s
}

func (a Gateway) GetStack() string {
	return a.Spec.Stack
}

func (a *Gateway) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

func (a Gateway) IsDebug() bool {
	return a.Spec.Debug
}

func (a Gateway) IsDev() bool {
	return a.Spec.Dev
}

//+kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
