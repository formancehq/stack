package v1beta2

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
