// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type GetWalletResponse struct {
	Data Wallet `json:"data"`
}

func (o *GetWalletResponse) GetData() Wallet {
	if o == nil {
		return Wallet{}
	}
	return o.Data
}
