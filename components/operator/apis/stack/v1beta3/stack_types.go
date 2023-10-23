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
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Ugly hack to be able to prevent secrets in testing
// Need to remove this
var ClientSecretGenerator = func() string {
	return uuid.NewString()
}

type IngressGlobalConfig struct {
	IngressConfig `json:",inline"`
	// +optional
	TLS *IngressTLS `json:"tls,omitempty"`
}

type DelegatedOIDCServerConfiguration struct {
	Issuer       string `json:"issuer,omitempty"`
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
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

func (cfg ClientConfiguration) WithAdditionalScopes(scopes ...string) ClientConfiguration {
	cfg.Scopes = append(cfg.Scopes, scopes...)
	return cfg
}

func (cfg ClientConfiguration) WithRedirectUris(redirectUris ...string) ClientConfiguration {
	cfg.RedirectUris = append(cfg.RedirectUris, redirectUris...)
	return cfg
}

func (cfg ClientConfiguration) WithPostLogoutRedirectUris(redirectUris ...string) ClientConfiguration {
	cfg.PostLogoutRedirectUris = append(cfg.PostLogoutRedirectUris, redirectUris...)
	return cfg
}

func NewClientConfiguration() ClientConfiguration {
	return ClientConfiguration{
		Scopes: []string{"openid"}, // Required scope
	}
}

type StaticClient struct {
	ClientConfiguration `json:",inline" yaml:",inline"`
	ID                  string `json:"id" yaml:"id"`
	// +optional
	Secrets []string `json:"secrets,omitempty" yaml:"secrets"`
}

type StackAuthSpec struct {
	DelegatedOIDCServer DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`
}

type StackStargateConfig struct {
	StargateServerURL string `json:"stargateServerURL"`
}

type StackServicesSpec struct {
	// +optional
	Ledger        StackServicePropertiesSpec `json:"ledger,omitempty"`
	Orchestration StackServicePropertiesSpec `json:"orchestration,omitempty"`
	Payments      StackServicePropertiesSpec `json:"payments,omitempty"`
	Wallets       StackServicePropertiesSpec `json:"wallets,omitempty"`
	Webhooks      StackServicePropertiesSpec `json:"webhooks,omitempty"`
	Control       StackServicePropertiesSpec `json:"control,omitempty"`
}

type StackServicePropertiesSpec struct {
	// +optional
	Disabled bool `json:"disabled"`
}

// StackSpec defines the desired state of Stack
type StackSpec struct {
	DevProperties `json:",inline"`
	Seed          string `json:"seed"`
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

func (in *StackServicesSpec) IsDisabled(service string) bool {
	valueOf := reflect.ValueOf(in).Elem().FieldByName(strcase.ToCamel(service))
	if valueOf.IsValid() {
		return valueOf.Interface().(StackServicePropertiesSpec).Disabled
	}
	return false
}

type ControlAuthentication struct {
	ClientID string
}

/* todo: group statuses by module */

type StackStatus struct {
	Status `json:",inline"`

	// +optional
	Ports map[string]map[string]int32 `json:"ports,omitempty"`

	// +optional
	StaticAuthClients map[string]StaticClient `json:"staticAuthClients,omitempty"`

	// +optional
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

func (s *Stack) GetScheme() string {
	if s.Spec.Scheme != "" {
		return s.Spec.Scheme
	}
	return "https"
}

func (s *Stack) URL() string {
	return fmt.Sprintf("%s://%s", s.GetScheme(), s.Spec.Host)
}

func (in *Stack) GetServiceNamespacedName(service string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Name,
		Name:      in.GetServiceName(service),
	}
}

func (in *Stack) GetServiceName(service string) string {
	return fmt.Sprintf("%s-%s", in.Name, service)
}

func (in *Stack) GetOrCreateClient(name string, configuration ClientConfiguration) StaticClient {
	if in.Status.StaticAuthClients == nil {
		in.Status.StaticAuthClients = map[string]StaticClient{}
	}
	if _, ok := in.Status.StaticAuthClients[name]; !ok {
		in.Status.StaticAuthClients[name] = StaticClient{
			ID:                  name,
			Secrets:             []string{ClientSecretGenerator()},
			ClientConfiguration: configuration,
		}
	}
	return in.Status.StaticAuthClients[name]
}

func (in *Stack) SetError(err error) {
	SetCondition(in, ConditionTypeError, metav1.ConditionTrue, err.Error())
	SetCondition(in, ConditionTypeReady, metav1.ConditionFalse)
}

func (in *Stack) SetReady() {
	RemoveCondition(in, ConditionTypeError)
	RemoveCondition(in, ConditionTypeProgressing)
	SetCondition(in, ConditionTypeReady, metav1.ConditionTrue)
}

func (in *Stack) SetProgressing() {
	SetCondition(in, ConditionTypeProgressing, metav1.ConditionTrue)
}

func (s *Stack) IsReady() bool {
	ret := collectionutils.Filter(s.Status.Conditions, func(t Condition) bool {
		return t.Type == ConditionTypeReady
	})
	if len(ret) == 0 {
		return false
	}
	return ret[0].Status == metav1.ConditionTrue
}

func (s *Stack) GetStaticClients(configuration *Configuration) []StaticClient {
	stackStaticClients := collectionutils.SliceFromMap(s.Status.StaticAuthClients)
	sort.SliceStable(
		stackStaticClients,
		func(i, j int) bool {
			return strings.Compare(stackStaticClients[i].ID, stackStaticClients[j].ID) < 0
		},
	)
	staticClients := append(configuration.Spec.Services.Auth.StaticClients, stackStaticClients...)
	staticClients = append(staticClients, s.Spec.Auth.StaticClients...)
	return staticClients
}

func (s *Stack) ModeChanged(configuration *Configuration) bool {
	return s.Status.LightMode != configuration.Spec.LightMode
}

func (s *Stack) IsDisabled(module string) bool {
	return s.Spec.Services.IsDisabled(module)
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
