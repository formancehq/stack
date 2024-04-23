// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type PaymentsAccountRaw struct {
}

type PaymentsAccount struct {
	AccountName  string    `json:"accountName"`
	ConnectorID  string    `json:"connectorID"`
	CreatedAt    time.Time `json:"createdAt"`
	DefaultAsset string    `json:"defaultAsset"`
	// Deprecated field: This will be removed in a future release, please migrate away from it as soon as possible.
	DefaultCurrency string              `json:"defaultCurrency"`
	ID              string              `json:"id"`
	Metadata        map[string]string   `json:"metadata"`
	Pools           []string            `json:"pools,omitempty"`
	Provider        *string             `json:"provider,omitempty"`
	Raw             *PaymentsAccountRaw `json:"raw"`
	Reference       string              `json:"reference"`
	Type            AccountType         `json:"type"`
}

func (p PaymentsAccount) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(p, "", false)
}

func (p *PaymentsAccount) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &p, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *PaymentsAccount) GetAccountName() string {
	if o == nil {
		return ""
	}
	return o.AccountName
}

func (o *PaymentsAccount) GetConnectorID() string {
	if o == nil {
		return ""
	}
	return o.ConnectorID
}

func (o *PaymentsAccount) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *PaymentsAccount) GetDefaultAsset() string {
	if o == nil {
		return ""
	}
	return o.DefaultAsset
}

func (o *PaymentsAccount) GetDefaultCurrency() string {
	if o == nil {
		return ""
	}
	return o.DefaultCurrency
}

func (o *PaymentsAccount) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *PaymentsAccount) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *PaymentsAccount) GetPools() []string {
	if o == nil {
		return nil
	}
	return o.Pools
}

func (o *PaymentsAccount) GetProvider() *string {
	if o == nil {
		return nil
	}
	return o.Provider
}

func (o *PaymentsAccount) GetRaw() *PaymentsAccountRaw {
	if o == nil {
		return nil
	}
	return o.Raw
}

func (o *PaymentsAccount) GetReference() string {
	if o == nil {
		return ""
	}
	return o.Reference
}

func (o *PaymentsAccount) GetType() AccountType {
	if o == nil {
		return AccountType("")
	}
	return o.Type
}
