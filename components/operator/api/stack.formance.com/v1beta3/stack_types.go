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

package v1beta3

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/stoewer/go-strcase"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
)

type IngressConfig struct {
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type resource struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourceProperties struct {
	// +optional
	Request *resource `json:"request,omitempty"`
	// +optional
	Limits *resource `json:"limits,omitempty"`
}

type IngressGlobalConfig struct {
	IngressConfig `json:",inline"`
	// +optional
	TLS *v1beta1.GatewayIngressTLS `json:"tls,omitempty"`
}

type ClientConfiguration struct {
	// +optional
	Public bool `json:"public"`
	// +optional
	Description *string `json:"description,omitempty"`
	// +optional
	RedirectUris []string `json:"redirectUris,omitempty" yaml:"redirectUris"`
	// +optional
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty" yaml:"PostLogoutRedirectUris"`
	// +optional
	Scopes []string `json:"scopes,omitempty"`
}

type StaticClient struct {
	ClientConfiguration `json:",inline" yaml:",inline"`
	ID                  string `json:"id" yaml:"id"`
	// +optional
	Secrets []string `json:"secrets,omitempty" yaml:"secrets"`
}

type StackAuthSpec struct {
	DelegatedOIDCServer v1beta1.DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`
	// +optional
	StaticClients []*v1beta1.AuthClientSpec `json:"staticClients,omitempty"`
}

type StackStargateConfig struct {
	StargateServerURL string `json:"stargateServerURL"`
}

type StackServicesSpec struct {
	// +optional
	Ledger         StackServicePropertiesSpec `json:"ledger,omitempty"`
	Orchestration  StackServicePropertiesSpec `json:"orchestration,omitempty"`
	Reconciliation StackServicePropertiesSpec `json:"reconciliation,omitempty"`
	Payments       StackServicePropertiesSpec `json:"payments,omitempty"`
	Wallets        StackServicePropertiesSpec `json:"wallets,omitempty"`
	Webhooks       StackServicePropertiesSpec `json:"webhooks,omitempty"`
	Control        StackServicePropertiesSpec `json:"control,omitempty"`
}

type StackServicePropertiesSpec struct {
	// +optional
	Disabled *bool `json:"disabled,omitempty"`
}

// StackSpec defines the desired state of Stack
type StackSpec struct {
	v1beta1.DevProperties `json:",inline"`
	Seed                  string `json:"seed"`
	// +kubebuilder:validation:Required
	Host string        `json:"host"`
	Auth StackAuthSpec `json:"auth"`

	// +optional
	Stargate *StackStargateConfig `json:"stargate,omitempty"`

	// +optional
	Versions string `json:"versions"`

	// +optional
	// +kubebuilder:default:="http"
	Scheme string `json:"scheme"`

	// +optional
	Disabled bool `json:"disabled"`

	// +optional
	Services StackServicesSpec `json:"services,omitempty"`
}

type ControlAuthentication struct {
	ClientID string
}

type StackStatus struct {
	Status `json:",inline"`

	//+optional
	Ready bool `json:"ready"`

	//+optional
	Ports map[string]map[string]int32 `json:"ports,omitempty"`

	//+optional
	StaticAuthClients map[string]StaticClient `json:"staticAuthClients,omitempty"`

	//+optional
	LightMode bool `json:"light"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Disable",type=string,JSONPath=".spec.disabled",description="Stack Disabled"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.versions",description="Stack Version"
//+kubebuilder:printcolumn:name="Configuration",type="string",JSONPath=".spec.seed",description="Stack Configuration"
//+kubebuilder:storageversion

// Stack is the Schema for the stacks API
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

func (*Stack) Hub() {}

func (in *StackServicesSpec) getDisabledProperty(service string) *bool {
	valueOf := reflect.ValueOf(in).Elem().FieldByName(strcase.UpperCamelCase(service))
	if valueOf.IsValid() {
		properties := valueOf.Interface().(StackServicePropertiesSpec)
		return properties.Disabled
	}
	return nil
}

func (in *StackServicesSpec) IsExplicitlyDisabled(service string) bool {
	disabled := in.getDisabledProperty(service)
	return disabled != nil && *disabled
}

func (in *StackServicesSpec) IsExplicitlyEnabled(service string) bool {
	disabled := in.getDisabledProperty(service)
	return disabled != nil && !*disabled
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
