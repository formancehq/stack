package v1beta3

import "github.com/formancehq/operator/api/formance.com/v1beta1"

type CommonServiceProperties struct {
	*v1beta1.DevProperties `json:",inline"`
	// +optional
	Disabled *bool `json:"disabled,omitempty"`
}
