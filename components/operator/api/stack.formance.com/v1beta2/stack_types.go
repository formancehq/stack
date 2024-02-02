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
	"encoding/json"

	authcomponentsv1beta2 "github.com/formancehq/operator/api/stack.formance.com/auth.components/v1beta2"
	"github.com/formancehq/operator/api/stack.formance.com/components/v1beta2"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

type IngressGlobalConfig struct {
	IngressConfig `json:",inline"`
	// +optional
	TLS *v1beta2.IngressTLS `json:"tls"`
}

type StackAuthSpec struct {
	DelegatedOIDCServer v1beta2.DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`
	// +optional
	StaticClients []authcomponentsv1beta2.StaticClient `json:"staticClients,omitempty"`
}

// StackSpec defines the desired state of Stack
type StackSpec struct {
	v1beta2.DevProperties `json:",inline"`
	Seed                  string        `json:"seed"`
	Host                  string        `json:"host"`
	Auth                  StackAuthSpec `json:"auth"`

	// +optional
	Versions string `json:"versions"`

	// +optional
	// +kubebuilder:default:="http"
	Scheme string `json:"scheme"`
}

type ControlAuthentication struct {
	ClientID string
}

type StackStatus struct {
	v1beta2.Status `json:",inline"`

	// +optional
	StaticAuthClients map[string]authcomponentsv1beta2.StaticClient `json:"staticAuthClients,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.progress`
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.versions",description="Stack Version"
//+kubebuilder:printcolumn:name="Configuration",type="string",JSONPath=".spec.seed",description="Stack Configuration"

// Stack is the Schema for the stacks API
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

func (stack *Stack) ConvertFrom(hubRaw conversion.Hub) error {
	hub := hubRaw.(*v1beta3.Stack)
	specAsRaw, err := json.Marshal(hub.Spec)
	if err != nil {
		return err
	}
	*stack = Stack{
		ObjectMeta: hub.ObjectMeta,
		TypeMeta:   hub.TypeMeta,
	}
	stack.APIVersion = GroupVersion.String()
	if err := json.Unmarshal(specAsRaw, &stack.Spec); err != nil {
		return err
	}
	return nil
}

func (stack *Stack) ConvertTo(hubRaw conversion.Hub) error {
	hub := hubRaw.(*v1beta3.Stack)
	specAsRaw, err := json.Marshal(stack.Spec)
	hub.ObjectMeta = stack.ObjectMeta
	hub.TypeMeta = stack.TypeMeta
	hub.TypeMeta.APIVersion = v1beta3.GroupVersion.String()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(specAsRaw, &hub.Spec); err != nil {
		return err
	}

	return nil
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
