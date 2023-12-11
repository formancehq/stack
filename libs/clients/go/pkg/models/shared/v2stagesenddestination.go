// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type V2StageSendDestination struct {
	Account *V2StageSendDestinationAccount `json:"account,omitempty"`
	Payment *V2StageSendDestinationPayment `json:"payment,omitempty"`
	Wallet  *V2StageSendDestinationWallet  `json:"wallet,omitempty"`
}

func (o *V2StageSendDestination) GetAccount() *V2StageSendDestinationAccount {
	if o == nil {
		return nil
	}
	return o.Account
}

func (o *V2StageSendDestination) GetPayment() *V2StageSendDestinationPayment {
	if o == nil {
		return nil
	}
	return o.Payment
}

func (o *V2StageSendDestination) GetWallet() *V2StageSendDestinationWallet {
	if o == nil {
		return nil
	}
	return o.Wallet
}
