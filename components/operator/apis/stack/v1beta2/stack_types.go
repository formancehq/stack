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
	"fmt"
	"reflect"
	"strings"

	authcomponentsv1beta2 "github.com/formancehq/operator/apis/auth.components/v1beta2"
	"github.com/formancehq/operator/apis/components/v1beta2"
	pkgapisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressGlobalConfig struct {
	IngressConfig `json:",inline"`
	// +optional
	TLS *pkgapisv1beta2.IngressTLS `json:"tls"`
}

type StackAuthSpec struct {
	DelegatedOIDCServer v1beta2.DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`
	// +optional
	StaticClients []authcomponentsv1beta2.StaticClient `json:"staticClients,omitempty"`
}

// StackSpec defines the desired state of Stack
type StackSpec struct {
	pkgapisv1beta2.DevProperties `json:",inline"`
	Seed                         string        `json:"seed"`
	Host                         string        `json:"host"`
	Auth                         StackAuthSpec `json:"auth"`

	// +optional
	Versions string `json:"versions"`

	// +optional
	// +kubebuilder:default:="http"
	Scheme string `json:"scheme"`
}

const (
	ConditionTypeStackNamespaceReady      = "NamespaceReady"
	ConditionTypeStackAuthReady           = "AuthReady"
	ConditionTypeStackLedgerReady         = "LedgerReady"
	ConditionTypeStackSearchReady         = "SearchReady"
	ConditionTypeStackControlReady        = "ControlReady"
	ConditionTypeStackPaymentsReady       = "PaymentsReady"
	ConditionTypeStackWebhooksReady       = "WebhooksReady"
	ConditionTypeStackOrchestrationReady  = "OrchestrationReady"
	ConditionTypeStackWalletsReady        = "WalletsReady"
	ConditionTypeStackCounterpartiesReady = "CounterpartiesReady"
	ConditionTypeStackMiddlewareReady     = "MiddlewareReady"
)

type ControlAuthentication struct {
	ClientID string
}

type StackStatus struct {
	pkgapisv1beta2.Status `json:",inline"`

	// +optional
	StaticAuthClients map[string]authcomponentsv1beta2.StaticClient `json:"staticAuthClients,omitempty"`
}

func (s *StackStatus) IsDirty(reference pkgapisv1beta2.Object) bool {
	if s.Status.IsDirty(reference) {
		return true
	}
	return !reflect.DeepEqual(reference.(*Stack).Status.StaticAuthClients, s.StaticAuthClients)
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.progress`
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

func NewStack(name string, spec StackSpec) Stack {
	return Stack{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: spec,
	}
}

func (*Stack) Hub() {}

func (s *Stack) GetScheme() string {
	if s.Spec.Scheme != "" {
		return s.Spec.Scheme
	}
	return "https"
}

func (s *Stack) URL() string {
	return fmt.Sprintf("%s://%s", s.GetScheme(), s.Spec.Host)
}

func (s *Stack) GetStatus() pkgapisv1beta2.Dirty {
	return &s.Status
}

func (s *Stack) IsDirty(t pkgapisv1beta2.Object) bool {
	return false
}

func (s *Stack) GetConditions() *pkgapisv1beta2.Conditions {
	return &s.Status.Conditions
}

func (s *Stack) SetNamespaceCreated() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackNamespaceReady, metav1.ConditionTrue)
}

func (s *Stack) SetNamespaceError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackNamespaceReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetAuthReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackAuthReady, metav1.ConditionTrue)
}

func (s *Stack) SetWebhooksReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackWebhooksReady, metav1.ConditionTrue)
}

func (s *Stack) SetAuthError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackAuthReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetLedgerReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackLedgerReady, metav1.ConditionTrue)
}

func (s *Stack) SetLedgerError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackLedgerReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetSearchReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackSearchReady, metav1.ConditionTrue)
}

func (s *Stack) SetSearchError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackSearchReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetControlReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackControlReady, metav1.ConditionTrue)
}

func (s *Stack) SetControlError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackControlReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetPaymentError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackPaymentsReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetWebhooksError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackWebhooksReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetOrchestrationError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackWebhooksReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetMiddlewareError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackMiddlewareReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetPaymentReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackPaymentsReady, metav1.ConditionTrue)
}

func (s *Stack) RemoveAuthStatus() {
	s.Status.RemoveCondition(ConditionTypeStackAuthReady)
}

func (s *Stack) RemoveSearchStatus() {
	s.Status.RemoveCondition(ConditionTypeStackSearchReady)
}

func (s *Stack) RemoveControlStatus() {
	s.Status.RemoveCondition(ConditionTypeStackControlReady)
}

func (in *Stack) RemovePaymentsStatus() {
	in.Status.RemoveCondition(ConditionTypeStackPaymentsReady)
}

func (in *Stack) RemoveWebhooksStatus() {
	in.Status.RemoveCondition(ConditionTypeStackWebhooksReady)
}

func (in *Stack) SetMiddlewareReady() {
	in.Status.RemoveCondition(ConditionTypeStackMiddlewareReady)
}

func (s *Stack) ServiceName(v string) string {
	return fmt.Sprintf("%s-%s", s.Name, strings.ToLower(v))
}

func (in *Stack) RemoveWalletsStatus() {
	in.Status.RemoveCondition(ConditionTypeStackWalletsReady)
}

func (s *Stack) SetWalletsError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackWalletsReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetWalletsReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackWalletsReady, metav1.ConditionTrue)
}

func (s *Stack) SetOrchestrationReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackOrchestrationReady, metav1.ConditionTrue)
}

func (in *Stack) RemoveCounterpartiesStatus() {
	in.Status.RemoveCondition(ConditionTypeStackCounterpartiesReady)
}

func (s *Stack) SetCounterpartiesError(msg string) {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackCounterpartiesReady, metav1.ConditionFalse, msg)
}

func (s *Stack) SetCounterpartiesReady() {
	pkgapisv1beta2.SetCondition(s, ConditionTypeStackCounterpartiesReady, metav1.ConditionTrue)
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
