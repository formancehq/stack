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

// WalletsPosting struct for WalletsPosting
type WalletsPosting struct {
	Amount      int64  `json:"amount"`
	Asset       string `json:"asset"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
}

// NewWalletsPosting instantiates a new WalletsPosting object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWalletsPosting(amount int64, asset string, destination string, source string) *WalletsPosting {
	this := WalletsPosting{}
	this.Amount = amount
	this.Asset = asset
	this.Destination = destination
	this.Source = source
	return &this
}

// NewWalletsPostingWithDefaults instantiates a new WalletsPosting object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWalletsPostingWithDefaults() *WalletsPosting {
	this := WalletsPosting{}
	return &this
}

// GetAmount returns the Amount field value
func (o *WalletsPosting) GetAmount() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Amount
}

// GetAmountOk returns a tuple with the Amount field value
// and a boolean to check if the value has been set.
func (o *WalletsPosting) GetAmountOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Amount, true
}

// SetAmount sets field value
func (o *WalletsPosting) SetAmount(v int64) {
	o.Amount = v
}

// GetAsset returns the Asset field value
func (o *WalletsPosting) GetAsset() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Asset
}

// GetAssetOk returns a tuple with the Asset field value
// and a boolean to check if the value has been set.
func (o *WalletsPosting) GetAssetOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Asset, true
}

// SetAsset sets field value
func (o *WalletsPosting) SetAsset(v string) {
	o.Asset = v
}

// GetDestination returns the Destination field value
func (o *WalletsPosting) GetDestination() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value
// and a boolean to check if the value has been set.
func (o *WalletsPosting) GetDestinationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Destination, true
}

// SetDestination sets field value
func (o *WalletsPosting) SetDestination(v string) {
	o.Destination = v
}

// GetSource returns the Source field value
func (o *WalletsPosting) GetSource() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *WalletsPosting) GetSourceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *WalletsPosting) SetSource(v string) {
	o.Source = v
}

func (o WalletsPosting) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["amount"] = o.Amount
	}
	if true {
		toSerialize["asset"] = o.Asset
	}
	if true {
		toSerialize["destination"] = o.Destination
	}
	if true {
		toSerialize["source"] = o.Source
	}
	return json.Marshal(toSerialize)
}

type NullableWalletsPosting struct {
	value *WalletsPosting
	isSet bool
}

func (v NullableWalletsPosting) Get() *WalletsPosting {
	return v.value
}

func (v *NullableWalletsPosting) Set(val *WalletsPosting) {
	v.value = val
	v.isSet = true
}

func (v NullableWalletsPosting) IsSet() bool {
	return v.isSet
}

func (v *NullableWalletsPosting) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWalletsPosting(val *WalletsPosting) *NullableWalletsPosting {
	return &NullableWalletsPosting{value: val, isSet: true}
}

func (v NullableWalletsPosting) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWalletsPosting) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
