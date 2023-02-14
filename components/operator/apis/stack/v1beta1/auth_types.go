package v1beta1

type DelegatedOIDCServerConfiguration struct {
	Issuer       string `json:"issuer,omitempty"`
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

type AuthSpec struct {
	ImageHolder `json:",inline"`
	// +optional
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	SigningKey string `json:"signingKey"`
	// +optional
	DelegatedOIDCServer *DelegatedOIDCServerConfiguration `json:"delegatedOIDCServer"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	Host string `json:"host,omitempty"`
	// +optional
	Scheme string `json:"scheme,omitempty"`
}

func (in *AuthSpec) GetScheme() string {
	if in.Scheme != "" {
		return in.Scheme
	}
	return "https"
}
