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

// checks if the StageSendSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StageSendSource{}

// StageSendSource struct for StageSendSource
type StageSendSource struct {
	Wallet *StageSendSourceWallet `json:"wallet,omitempty"`
	Account *StageSendSourceAccount `json:"account,omitempty"`
	Payment *StageSendSourcePayment `json:"payment,omitempty"`
}

// NewStageSendSource instantiates a new StageSendSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStageSendSource() *StageSendSource {
	this := StageSendSource{}
	return &this
}

// NewStageSendSourceWithDefaults instantiates a new StageSendSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStageSendSourceWithDefaults() *StageSendSource {
	this := StageSendSource{}
	return &this
}

// GetWallet returns the Wallet field value if set, zero value otherwise.
func (o *StageSendSource) GetWallet() StageSendSourceWallet {
	if o == nil || IsNil(o.Wallet) {
		var ret StageSendSourceWallet
		return ret
	}
	return *o.Wallet
}

// GetWalletOk returns a tuple with the Wallet field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendSource) GetWalletOk() (*StageSendSourceWallet, bool) {
	if o == nil || IsNil(o.Wallet) {
		return nil, false
	}
	return o.Wallet, true
}

// HasWallet returns a boolean if a field has been set.
func (o *StageSendSource) HasWallet() bool {
	if o != nil && !IsNil(o.Wallet) {
		return true
	}

	return false
}

// SetWallet gets a reference to the given StageSendSourceWallet and assigns it to the Wallet field.
func (o *StageSendSource) SetWallet(v StageSendSourceWallet) {
	o.Wallet = &v
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *StageSendSource) GetAccount() StageSendSourceAccount {
	if o == nil || IsNil(o.Account) {
		var ret StageSendSourceAccount
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendSource) GetAccountOk() (*StageSendSourceAccount, bool) {
	if o == nil || IsNil(o.Account) {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *StageSendSource) HasAccount() bool {
	if o != nil && !IsNil(o.Account) {
		return true
	}

	return false
}

// SetAccount gets a reference to the given StageSendSourceAccount and assigns it to the Account field.
func (o *StageSendSource) SetAccount(v StageSendSourceAccount) {
	o.Account = &v
}

// GetPayment returns the Payment field value if set, zero value otherwise.
func (o *StageSendSource) GetPayment() StageSendSourcePayment {
	if o == nil || IsNil(o.Payment) {
		var ret StageSendSourcePayment
		return ret
	}
	return *o.Payment
}

// GetPaymentOk returns a tuple with the Payment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSendSource) GetPaymentOk() (*StageSendSourcePayment, bool) {
	if o == nil || IsNil(o.Payment) {
		return nil, false
	}
	return o.Payment, true
}

// HasPayment returns a boolean if a field has been set.
func (o *StageSendSource) HasPayment() bool {
	if o != nil && !IsNil(o.Payment) {
		return true
	}

	return false
}

// SetPayment gets a reference to the given StageSendSourcePayment and assigns it to the Payment field.
func (o *StageSendSource) SetPayment(v StageSendSourcePayment) {
	o.Payment = &v
}

func (o StageSendSource) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StageSendSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Wallet) {
		toSerialize["wallet"] = o.Wallet
	}
	if !IsNil(o.Account) {
		toSerialize["account"] = o.Account
	}
	if !IsNil(o.Payment) {
		toSerialize["payment"] = o.Payment
	}
	return toSerialize, nil
}

type NullableStageSendSource struct {
	value *StageSendSource
	isSet bool
}

func (v NullableStageSendSource) Get() *StageSendSource {
	return v.value
}

func (v *NullableStageSendSource) Set(val *StageSendSource) {
	v.value = val
	v.isSet = true
}

func (v NullableStageSendSource) IsSet() bool {
	return v.isSet
}

func (v *NullableStageSendSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStageSendSource(val *StageSendSource) *NullableStageSendSource {
	return &NullableStageSendSource{value: val, isSet: true}
}

func (v NullableStageSendSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStageSendSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


