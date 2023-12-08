// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"math/big"
)

type AggregateBalancesResponse struct {
	Data map[string]*big.Int `json:"data"`
}

func (a AggregateBalancesResponse) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AggregateBalancesResponse) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AggregateBalancesResponse) GetData() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.Data
}
