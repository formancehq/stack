// +kubebuilder:object:generate=true
package v1beta3

import (
	"fmt"
)

type PostgresConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	// +optional
	Username string `json:"username"`
	// +optional
	Password string `json:"password"`
	// +optional
	CredentialsFromSecret string `json:"credentialsFromSecret"`
	// +optional
	DisableSSLMode bool `json:"disableSSLMode"`
}

func (in *PostgresConfig) DSN() string {
	queryParams := ""
	if in.DisableSSLMode {
		queryParams = "?sslmode=disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/postgres%s", in.Username, in.Password, in.Host, in.Port, queryParams)
}
