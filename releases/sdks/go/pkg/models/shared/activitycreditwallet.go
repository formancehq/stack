// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type ActivityCreditWallet struct {
	Data *CreditWalletRequest `json:"data,omitempty"`
	ID   *string              `json:"id,omitempty"`
}

func (o *ActivityCreditWallet) GetData() *CreditWalletRequest {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *ActivityCreditWallet) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}
