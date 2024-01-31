// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
)

type StripeConfig struct {
	APIKey string `json:"apiKey"`
	Name   string `json:"name"`
	// Number of BalanceTransaction to fetch at each polling interval.
	//
	PageSize *int64 `default:"10" json:"pageSize"`
	// The frequency at which the connector will try to fetch new BalanceTransaction objects from Stripe API.
	//
	PollingPeriod *string `default:"120s" json:"pollingPeriod"`
}

func (s StripeConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(s, "", false)
}

func (s *StripeConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &s, "", false, true); err != nil {
		return err
	}
	return nil
}

func (o *StripeConfig) GetAPIKey() string {
	if o == nil {
		return ""
	}
	return o.APIKey
}

func (o *StripeConfig) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *StripeConfig) GetPageSize() *int64 {
	if o == nil {
		return nil
	}
	return o.PageSize
}

func (o *StripeConfig) GetPollingPeriod() *string {
	if o == nil {
		return nil
	}
	return o.PollingPeriod
}