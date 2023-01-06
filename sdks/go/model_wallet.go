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

// Wallet struct for Wallet
type Wallet struct {
	// The unique ID of the wallet.
	Id string `json:"id"`
	// Metadata associated with the wallet.
	Metadata map[string]interface{} `json:"metadata"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewWallet instantiates a new Wallet object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWallet(id string, metadata map[string]interface{}, name string, createdAt time.Time) *Wallet {
	this := Wallet{}
	this.Id = id
	this.Metadata = metadata
	this.Name = name
	this.CreatedAt = createdAt
	return &this
}

// NewWalletWithDefaults instantiates a new Wallet object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWalletWithDefaults() *Wallet {
	this := Wallet{}
	return &this
}

// GetId returns the Id field value
func (o *Wallet) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Wallet) GetIdOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Wallet) SetId(v string) {
	o.Id = v
}

// GetMetadata returns the Metadata field value
func (o *Wallet) GetMetadata() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *Wallet) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil {
    return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// SetMetadata sets field value
func (o *Wallet) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetName returns the Name field value
func (o *Wallet) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Wallet) GetNameOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Wallet) SetName(v string) {
	o.Name = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Wallet) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Wallet) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
    return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Wallet) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

func (o Wallet) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["metadata"] = o.Metadata
	}
	if true {
		toSerialize["name"] = o.Name
	}
	if true {
		toSerialize["createdAt"] = o.CreatedAt
	}
	return json.Marshal(toSerialize)
}

type NullableWallet struct {
	value *Wallet
	isSet bool
}

func (v NullableWallet) Get() *Wallet {
	return v.value
}

func (v *NullableWallet) Set(val *Wallet) {
	v.value = val
	v.isSet = true
}

func (v NullableWallet) IsSet() bool {
	return v.isSet
}

func (v *NullableWallet) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWallet(val *Wallet) *NullableWallet {
	return &NullableWallet{value: val, isSet: true}
}

func (v NullableWallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWallet) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
