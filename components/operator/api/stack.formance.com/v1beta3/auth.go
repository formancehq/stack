package v1beta3

type AuthConfig struct {
	// +optional
	ReadKeySetMaxRetries int `json:"readKeySetMaxRetries"`
	// +optional
	CheckScopes bool `json:"checkScopes"`
}
