// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type V2DebitWalletRequest struct {
	Amount      V2Monetary `json:"amount"`
	Balances    []string   `json:"balances,omitempty"`
	Description *string    `json:"description,omitempty"`
	Destination *V2Subject `json:"destination,omitempty"`
	// Metadata associated with the wallet.
	Metadata map[string]string `json:"metadata"`
	// Set to true to create a pending hold. If false, the wallet will be debited immediately.
	Pending *bool `json:"pending,omitempty"`
	// cannot be used in conjunction with `pending` property
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

func (v V2DebitWalletRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2DebitWalletRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2DebitWalletRequest) GetAmount() V2Monetary {
	if o == nil {
		return V2Monetary{}
	}
	return o.Amount
}

func (o *V2DebitWalletRequest) GetBalances() []string {
	if o == nil {
		return nil
	}
	return o.Balances
}

func (o *V2DebitWalletRequest) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *V2DebitWalletRequest) GetDestination() *V2Subject {
	if o == nil {
		return nil
	}
	return o.Destination
}

func (o *V2DebitWalletRequest) GetDestinationAccount() *V2LedgerAccountSubject {
	if v := o.GetDestination(); v != nil {
		return v.V2LedgerAccountSubject
	}
	return nil
}

func (o *V2DebitWalletRequest) GetDestinationWallet() *V2WalletSubject {
	if v := o.GetDestination(); v != nil {
		return v.V2WalletSubject
	}
	return nil
}

func (o *V2DebitWalletRequest) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *V2DebitWalletRequest) GetPending() *bool {
	if o == nil {
		return nil
	}
	return o.Pending
}

func (o *V2DebitWalletRequest) GetTimestamp() *time.Time {
	if o == nil {
		return nil
	}
	return o.Timestamp
}
