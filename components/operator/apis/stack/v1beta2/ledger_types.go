package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
)

// +kubebuilder:object:generate=true
type LedgerSpec struct {
	componentsv1beta2.Scalable `json:",inline"`
	Postgres                   componentsv1beta2.PostgresConfig `json:"postgres"`
	// +optional
	LockingStrategy componentsv1beta2.LockingStrategy `json:"locking"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
}
