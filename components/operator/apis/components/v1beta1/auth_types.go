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
	"github.com/formancehq/operator/apis/auth.components/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OAuth2ConfigSpec struct {

	// +required
	IntrospectUrl string `json:"introspectUrl"`

	// +optional
	Audiences []string `json:"audiences"`

	// +optional
	AudienceWildcard bool `json:"audienceWildcard"`

	//+optional
	ProtectedByScopes bool `json:"ProtectedByScopes"`
}

type HTTPBasicConfigSpec struct {
	// +optional
	Enabled bool `json:"enabled"`

	// +optional
	Credentials map[string]string `json:"credentials"`
}

type AuthConfigSpec struct {
	// +optional
	OAuth2 *OAuth2ConfigSpec `json:"oauth2,omitempty"`

	// +optional
	HTTPBasic *HTTPBasicConfigSpec `json:"basic,omitempty"`
}

type DelegatedOIDCServerConfiguration struct {
	// +optional
	Issuer string `json:"issuer,omitempty"`
	// +optional
	IssuerFrom *ConfigSource `json:"issuerFrom,omitempty"`
	// +optional
	ClientID string `json:"clientID,omitempty"`
	// +optional
	ClientIDFrom *ConfigSource `json:"clientIDFrom,omitempty"`
	// +optional
	ClientSecret string `json:"clientSecret,omitempty"`
	// +optional
	ClientSecretFrom *ConfigSource `json:"clientSecretFrom,omitempty"`
}

// AuthSpec defines the desired state of Auth
type AuthSpec struct {
	Scalable    `json:",inline"`
	ImageHolder `json:",inline"`
	Postgres    PostgresConfigCreateDatabase `json:"postgres"`
	BaseURL     string                       `json:"baseURL"`

	// SigningKey is a private key
	// The signing key is used by the server to sign JWT tokens
	// The value of this config will be copied to a secret and injected inside
	// the env vars of the server using secret mapping.
	// If not specified, a key will be automatically generated.
	// +optional
	SigningKey string `json:"signingKey"`
	DevMode    bool   `json:"devMode"`
	// +optional
	Ingress *IngressSpec `json:"ingress"`

	DelegatedOIDCServer DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`

	// +optional
	Monitoring *MonitoringSpec `json:"monitoring"`

	// +optional
	StaticClients []v1beta1.StaticClient `json:"staticClients"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector

// Auth is the Schema for the auths API
type Auth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthSpec          `json:"spec,omitempty"`
	Status ReplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AuthList contains a list of Auth
type AuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Auth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Auth{}, &AuthList{})
}
