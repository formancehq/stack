// +kubebuilder:object:generate=true
package v1beta3

type PostgresConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	// +optional
	Username string `json:"username"`
	// +optional
	Password string `json:"password"`

	// +optional
	Debug bool `json:"debug"`
	// +optional
	CredentialsFromSecret string `json:"credentialsFromSecret"`
	// +optional
	DisableSSLMode bool `json:"disableSSLMode"`
}
