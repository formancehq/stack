/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the StageSendDestination type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StageSendDestination{}

// StageSendDestination struct for StageSendDestination
type StageSendDestination struct {
	Wallet *StageSendSourceWallet `json:"wallet,omitempty"`
	Account *StageSendSourceAccount `json:"account,omitempty"`
	Payment *StageSendDestinationPayment `json:"payment,omitempty"`
}

// NewStageSendDestination instantiates a new StageSendDestination object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStageSendDestination() *StageSendDestination {
	this := StageSendDestination{}
	return &this
}

// NewStageSendDestinationWithDefaults instantiates a new StageSendDestination object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStageSendDestinationWithDefaults() *StageSendDestination {
	this := StageSendDestination{}
	return &this
}

// GetWallet returns the Wallet field value if set, zero value otherwise.
func (o *StageSendDestination) GetWallet() StageSendSourceWallet {
	if o == nil || isNil(o.Wallet) {
		var ret StageSendSourceWallet
		return ret
	}
	return *o.Wallet
}

// GetWalletOk returns a tuple with the Wallet field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendDestination) GetWalletOk() (*StageSendSourceWallet, bool) {
	if o == nil || isNil(o.Wallet) {
		return nil, false
	}
	return o.Wallet, true
}

// HasWallet returns a boolean if a field has been set.
func (o *StageSendDestination) HasWallet() bool {
	if o != nil && !isNil(o.Wallet) {
		return true
	}

	return false
}

// SetWallet gets a reference to the given StageSendSourceWallet and assigns it to the Wallet field.
func (o *StageSendDestination) SetWallet(v StageSendSourceWallet) {
	o.Wallet = &v
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *StageSendDestination) GetAccount() StageSendSourceAccount {
	if o == nil || isNil(o.Account) {
		var ret StageSendSourceAccount
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendDestination) GetAccountOk() (*StageSendSourceAccount, bool) {
	if o == nil || isNil(o.Account) {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *StageSendDestination) HasAccount() bool {
	if o != nil && !isNil(o.Account) {
		return true
	}

	return false
}

// SetAccount gets a reference to the given StageSendSourceAccount and assigns it to the Account field.
func (o *StageSendDestination) SetAccount(v StageSendSourceAccount) {
	o.Account = &v
}

// GetPayment returns the Payment field value if set, zero value otherwise.
func (o *StageSendDestination) GetPayment() StageSendDestinationPayment {
	if o == nil || isNil(o.Payment) {
		var ret StageSendDestinationPayment
		return ret
	}
	return *o.Payment
}

// GetPaymentOk returns a tuple with the Payment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendDestination) GetPaymentOk() (*StageSendDestinationPayment, bool) {
	if o == nil || isNil(o.Payment) {
		return nil, false
	}
	return o.Payment, true
}

// HasPayment returns a boolean if a field has been set.
func (o *StageSendDestination) HasPayment() bool {
	if o != nil && !isNil(o.Payment) {
		return true
	}

	return false
}

// SetPayment gets a reference to the given StageSendDestinationPayment and assigns it to the Payment field.
func (o *StageSendDestination) SetPayment(v StageSendDestinationPayment) {
	o.Payment = &v
}

func (o StageSendDestination) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StageSendDestination) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Wallet) {
		toSerialize["wallet"] = o.Wallet
	}
	if !isNil(o.Account) {
		toSerialize["account"] = o.Account
	}
	if !isNil(o.Payment) {
		toSerialize["payment"] = o.Payment
	}
	return toSerialize, nil
}

type NullableStageSendDestination struct {
	value *StageSendDestination
	isSet bool
}

func (v NullableStageSendDestination) Get() *StageSendDestination {
	return v.value
}

func (v *NullableStageSendDestination) Set(val *StageSendDestination) {
	v.value = val
	v.isSet = true
}

func (v NullableStageSendDestination) IsSet() bool {
	return v.isSet
}

func (v *NullableStageSendDestination) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStageSendDestination(val *StageSendDestination) *NullableStageSendDestination {
	return &NullableStageSendDestination{value: val, isSet: true}
}

func (v NullableStageSendDestination) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStageSendDestination) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


