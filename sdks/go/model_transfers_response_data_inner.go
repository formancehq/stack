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

// checks if the TransfersResponseDataInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TransfersResponseDataInner{}

// TransfersResponseDataInner struct for TransfersResponseDataInner
type TransfersResponseDataInner struct {
	Id *string `json:"id,omitempty"`
	Amount *int64 `json:"amount,omitempty"`
	Asset *string `json:"asset,omitempty"`
	Destination *string `json:"destination,omitempty"`
	Source *string `json:"source,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Status *string `json:"status,omitempty"`
	Error *string `json:"error,omitempty"`
}

// NewTransfersResponseDataInner instantiates a new TransfersResponseDataInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTransfersResponseDataInner() *TransfersResponseDataInner {
	this := TransfersResponseDataInner{}
	return &this
}

// NewTransfersResponseDataInnerWithDefaults instantiates a new TransfersResponseDataInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTransfersResponseDataInnerWithDefaults() *TransfersResponseDataInner {
	this := TransfersResponseDataInner{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetId() string {
	if o == nil || isNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetIdOk() (*string, bool) {
	if o == nil || isNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasId() bool {
	if o != nil && !isNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *TransfersResponseDataInner) SetId(v string) {
	o.Id = &v
}

// GetAmount returns the Amount field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetAmount() int64 {
	if o == nil || isNil(o.Amount) {
		var ret int64
		return ret
	}
	return *o.Amount
}

// GetAmountOk returns a tuple with the Amount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetAmountOk() (*int64, bool) {
	if o == nil || isNil(o.Amount) {
		return nil, false
	}
	return o.Amount, true
}

// HasAmount returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasAmount() bool {
	if o != nil && !isNil(o.Amount) {
		return true
	}

	return false
}

// SetAmount gets a reference to the given int64 and assigns it to the Amount field.
func (o *TransfersResponseDataInner) SetAmount(v int64) {
	o.Amount = &v
}

// GetAsset returns the Asset field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetAsset() string {
	if o == nil || isNil(o.Asset) {
		var ret string
		return ret
	}
	return *o.Asset
}

// GetAssetOk returns a tuple with the Asset field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetAssetOk() (*string, bool) {
	if o == nil || isNil(o.Asset) {
		return nil, false
	}
	return o.Asset, true
}

// HasAsset returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasAsset() bool {
	if o != nil && !isNil(o.Asset) {
		return true
	}

	return false
}

// SetAsset gets a reference to the given string and assigns it to the Asset field.
func (o *TransfersResponseDataInner) SetAsset(v string) {
	o.Asset = &v
}

// GetDestination returns the Destination field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetDestination() string {
	if o == nil || isNil(o.Destination) {
		var ret string
		return ret
	}
	return *o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetDestinationOk() (*string, bool) {
	if o == nil || isNil(o.Destination) {
		return nil, false
	}
	return o.Destination, true
}

// HasDestination returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasDestination() bool {
	if o != nil && !isNil(o.Destination) {
		return true
	}

	return false
}

// SetDestination gets a reference to the given string and assigns it to the Destination field.
func (o *TransfersResponseDataInner) SetDestination(v string) {
	o.Destination = &v
}

// GetSource returns the Source field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetSource() string {
	if o == nil || isNil(o.Source) {
		var ret string
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetSourceOk() (*string, bool) {
	if o == nil || isNil(o.Source) {
		return nil, false
	}
	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasSource() bool {
	if o != nil && !isNil(o.Source) {
		return true
	}

	return false
}

// SetSource gets a reference to the given string and assigns it to the Source field.
func (o *TransfersResponseDataInner) SetSource(v string) {
	o.Source = &v
}

// GetCurrency returns the Currency field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetCurrency() string {
	if o == nil || isNil(o.Currency) {
		var ret string
		return ret
	}
	return *o.Currency
}

// GetCurrencyOk returns a tuple with the Currency field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetCurrencyOk() (*string, bool) {
	if o == nil || isNil(o.Currency) {
		return nil, false
	}
	return o.Currency, true
}

// HasCurrency returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasCurrency() bool {
	if o != nil && !isNil(o.Currency) {
		return true
	}

	return false
}

// SetCurrency gets a reference to the given string and assigns it to the Currency field.
func (o *TransfersResponseDataInner) SetCurrency(v string) {
	o.Currency = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetStatus() string {
	if o == nil || isNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetStatusOk() (*string, bool) {
	if o == nil || isNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasStatus() bool {
	if o != nil && !isNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *TransfersResponseDataInner) SetStatus(v string) {
	o.Status = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *TransfersResponseDataInner) GetError() string {
	if o == nil || isNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponseDataInner) GetErrorOk() (*string, bool) {
	if o == nil || isNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *TransfersResponseDataInner) HasError() bool {
	if o != nil && !isNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *TransfersResponseDataInner) SetError(v string) {
	o.Error = &v
}

func (o TransfersResponseDataInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TransfersResponseDataInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !isNil(o.Amount) {
		toSerialize["amount"] = o.Amount
	}
	if !isNil(o.Asset) {
		toSerialize["asset"] = o.Asset
	}
	if !isNil(o.Destination) {
		toSerialize["destination"] = o.Destination
	}
	if !isNil(o.Source) {
		toSerialize["source"] = o.Source
	}
	if !isNil(o.Currency) {
		toSerialize["currency"] = o.Currency
	}
	if !isNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !isNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	return toSerialize, nil
}

type NullableTransfersResponseDataInner struct {
	value *TransfersResponseDataInner
	isSet bool
}

func (v NullableTransfersResponseDataInner) Get() *TransfersResponseDataInner {
	return v.value
}

func (v *NullableTransfersResponseDataInner) Set(val *TransfersResponseDataInner) {
	v.value = val
	v.isSet = true
}

func (v NullableTransfersResponseDataInner) IsSet() bool {
	return v.isSet
}

func (v *NullableTransfersResponseDataInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTransfersResponseDataInner(val *TransfersResponseDataInner) *NullableTransfersResponseDataInner {
	return &NullableTransfersResponseDataInner{value: val, isSet: true}
}

func (v NullableTransfersResponseDataInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTransfersResponseDataInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


