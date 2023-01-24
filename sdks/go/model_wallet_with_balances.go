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
	"time"
)

// checks if the WalletWithBalances type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WalletWithBalances{}

// WalletWithBalances struct for WalletWithBalances
type WalletWithBalances struct {
	// The unique ID of the wallet.
	Id string `json:"id"`
	// Metadata associated with the wallet.
	Metadata map[string]interface{} `json:"metadata"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Balances WalletWithBalancesBalances `json:"balances"`
	Ledger string `json:"ledger"`
}

// NewWalletWithBalances instantiates a new WalletWithBalances object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWalletWithBalances(id string, metadata map[string]interface{}, name string, createdAt time.Time, balances WalletWithBalancesBalances, ledger string) *WalletWithBalances {
	this := WalletWithBalances{}
	this.Id = id
	this.Metadata = metadata
	this.Name = name
	this.CreatedAt = createdAt
	this.Balances = balances
	this.Ledger = ledger
	return &this
}

// NewWalletWithBalancesWithDefaults instantiates a new WalletWithBalances object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWalletWithBalancesWithDefaults() *WalletWithBalances {
	this := WalletWithBalances{}
	return &this
}

// GetId returns the Id field value
func (o *WalletWithBalances) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WalletWithBalances) SetId(v string) {
	o.Id = v
}

// GetMetadata returns the Metadata field value
func (o *WalletWithBalances) GetMetadata() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// SetMetadata sets field value
func (o *WalletWithBalances) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetName returns the Name field value
func (o *WalletWithBalances) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WalletWithBalances) SetName(v string) {
	o.Name = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *WalletWithBalances) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *WalletWithBalances) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetBalances returns the Balances field value
func (o *WalletWithBalances) GetBalances() WalletWithBalancesBalances {
	if o == nil {
		var ret WalletWithBalancesBalances
		return ret
	}

	return o.Balances
}

// GetBalancesOk returns a tuple with the Balances field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetBalancesOk() (*WalletWithBalancesBalances, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Balances, true
}

// SetBalances sets field value
func (o *WalletWithBalances) SetBalances(v WalletWithBalancesBalances) {
	o.Balances = v
}

// GetLedger returns the Ledger field value
func (o *WalletWithBalances) GetLedger() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Ledger
}

// GetLedgerOk returns a tuple with the Ledger field value
// and a boolean to check if the value has been set.
func (o *WalletWithBalances) GetLedgerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ledger, true
}

// SetLedger sets field value
func (o *WalletWithBalances) SetLedger(v string) {
	o.Ledger = v
}

func (o WalletWithBalances) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WalletWithBalances) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["metadata"] = o.Metadata
	toSerialize["name"] = o.Name
	toSerialize["createdAt"] = o.CreatedAt
	toSerialize["balances"] = o.Balances
	toSerialize["ledger"] = o.Ledger
	return toSerialize, nil
}

type NullableWalletWithBalances struct {
	value *WalletWithBalances
	isSet bool
}

func (v NullableWalletWithBalances) Get() *WalletWithBalances {
	return v.value
}

func (v *NullableWalletWithBalances) Set(val *WalletWithBalances) {
	v.value = val
	v.isSet = true
}

func (v NullableWalletWithBalances) IsSet() bool {
	return v.isSet
}

func (v *NullableWalletWithBalances) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWalletWithBalances(val *WalletWithBalances) *NullableWalletWithBalances {
	return &NullableWalletWithBalances{value: val, isSet: true}
}

func (v NullableWalletWithBalances) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWalletWithBalances) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


