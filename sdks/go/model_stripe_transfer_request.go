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

// checks if the StripeTransferRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StripeTransferRequest{}

// StripeTransferRequest struct for StripeTransferRequest
type StripeTransferRequest struct {
	Amount *int64 `json:"amount,omitempty"`
	Asset *string `json:"asset,omitempty"`
	Destination *string `json:"destination,omitempty"`
	// A set of key/value pairs that you can attach to a transfer object. It can be useful for storing additional information about the transfer in a structured format. 
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewStripeTransferRequest instantiates a new StripeTransferRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStripeTransferRequest() *StripeTransferRequest {
	this := StripeTransferRequest{}
	return &this
}

// NewStripeTransferRequestWithDefaults instantiates a new StripeTransferRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStripeTransferRequestWithDefaults() *StripeTransferRequest {
	this := StripeTransferRequest{}
	return &this
}

// GetAmount returns the Amount field value if set, zero value otherwise.
func (o *StripeTransferRequest) GetAmount() int64 {
	if o == nil || IsNil(o.Amount) {
		var ret int64
		return ret
	}
	return *o.Amount
}

// GetAmountOk returns a tuple with the Amount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeTransferRequest) GetAmountOk() (*int64, bool) {
	if o == nil || IsNil(o.Amount) {
		return nil, false
	}
	return o.Amount, true
}

// HasAmount returns a boolean if a field has been set.
func (o *StripeTransferRequest) HasAmount() bool {
	if o != nil && !IsNil(o.Amount) {
		return true
	}

	return false
}

// SetAmount gets a reference to the given int64 and assigns it to the Amount field.
func (o *StripeTransferRequest) SetAmount(v int64) {
	o.Amount = &v
}

// GetAsset returns the Asset field value if set, zero value otherwise.
func (o *StripeTransferRequest) GetAsset() string {
	if o == nil || IsNil(o.Asset) {
		var ret string
		return ret
	}
	return *o.Asset
}

// GetAssetOk returns a tuple with the Asset field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeTransferRequest) GetAssetOk() (*string, bool) {
	if o == nil || IsNil(o.Asset) {
		return nil, false
	}
	return o.Asset, true
}

// HasAsset returns a boolean if a field has been set.
func (o *StripeTransferRequest) HasAsset() bool {
	if o != nil && !IsNil(o.Asset) {
		return true
	}

	return false
}

// SetAsset gets a reference to the given string and assigns it to the Asset field.
func (o *StripeTransferRequest) SetAsset(v string) {
	o.Asset = &v
}

// GetDestination returns the Destination field value if set, zero value otherwise.
func (o *StripeTransferRequest) GetDestination() string {
	if o == nil || IsNil(o.Destination) {
		var ret string
		return ret
	}
	return *o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeTransferRequest) GetDestinationOk() (*string, bool) {
	if o == nil || IsNil(o.Destination) {
		return nil, false
	}
	return o.Destination, true
}

// HasDestination returns a boolean if a field has been set.
func (o *StripeTransferRequest) HasDestination() bool {
	if o != nil && !IsNil(o.Destination) {
		return true
	}

	return false
}

// SetDestination gets a reference to the given string and assigns it to the Destination field.
func (o *StripeTransferRequest) SetDestination(v string) {
	o.Destination = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *StripeTransferRequest) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeTransferRequest) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *StripeTransferRequest) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *StripeTransferRequest) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

func (o StripeTransferRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StripeTransferRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Amount) {
		toSerialize["amount"] = o.Amount
	}
	if !IsNil(o.Asset) {
		toSerialize["asset"] = o.Asset
	}
	if !IsNil(o.Destination) {
		toSerialize["destination"] = o.Destination
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullableStripeTransferRequest struct {
	value *StripeTransferRequest
	isSet bool
}

func (v NullableStripeTransferRequest) Get() *StripeTransferRequest {
	return v.value
}

func (v *NullableStripeTransferRequest) Set(val *StripeTransferRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableStripeTransferRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableStripeTransferRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStripeTransferRequest(val *StripeTransferRequest) *NullableStripeTransferRequest {
	return &NullableStripeTransferRequest{value: val, isSet: true}
}

func (v NullableStripeTransferRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStripeTransferRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


