// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type PoolBalances struct {
	Balances []PoolBalance `json:"balances"`
}

func (o *PoolBalances) GetBalances() []PoolBalance {
	if o == nil {
		return []PoolBalance{}
	}
	return o.Balances
}
