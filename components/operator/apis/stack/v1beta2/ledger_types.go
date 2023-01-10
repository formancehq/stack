package v1beta2

import (
	authcomponentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	"github.com/formancehq/operator/pkg/apis/v1beta2"
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/typeutils"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=true
type LedgerSpec struct {
	apisv1beta2.Scalable `json:",inline"`
	Postgres             apisv1beta2.PostgresConfig `json:"postgres"`
	// +optional
	LockingStrategy authcomponentsv1beta2.LockingStrategy `json:"locking"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
}

func (l LedgerSpec) DatabaseSpec() apisv1beta2.PostgresConfigWithDatabase {
	return apisv1beta2.PostgresConfigWithDatabase{
		PostgresConfig: l.Postgres,
		Database:       "",
	}
}

func (in *LedgerSpec) Validate() field.ErrorList {
	if in == nil {
		return nil
	}
	ret := typeutils.Map(in.Postgres.Validate(), v1beta2.AddPrefixToFieldError("postgres"))
	ret = append(ret, typeutils.Map(in.LockingStrategy.Validate(), v1beta2.AddPrefixToFieldError("locking"))...)
	return ret
}
