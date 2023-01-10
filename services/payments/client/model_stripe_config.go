/*
Payments API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the StripeConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StripeConfig{}

// StripeConfig struct for StripeConfig
type StripeConfig struct {
	// The frequency at which the connector will try to fetch new BalanceTransaction objects from Stripe api
	PollingPeriod *string `json:"pollingPeriod,omitempty"`
	ApiKey        string  `json:"apiKey"`
	// Number of BalanceTransaction to fetch at each polling interval.
	PageSize *float32 `json:"pageSize,omitempty"`
}

// NewStripeConfig instantiates a new StripeConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStripeConfig(apiKey string) *StripeConfig {
	this := StripeConfig{}
	var pollingPeriod string = "120s"
	this.PollingPeriod = &pollingPeriod
	this.ApiKey = apiKey
	var pageSize float32 = 10
	this.PageSize = &pageSize
	return &this
}

// NewStripeConfigWithDefaults instantiates a new StripeConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStripeConfigWithDefaults() *StripeConfig {
	this := StripeConfig{}
	var pollingPeriod string = "120s"
	this.PollingPeriod = &pollingPeriod
	var pageSize float32 = 10
	this.PageSize = &pageSize
	return &this
}

// GetPollingPeriod returns the PollingPeriod field value if set, zero value otherwise.
func (o *StripeConfig) GetPollingPeriod() string {
	if o == nil || isNil(o.PollingPeriod) {
		var ret string
		return ret
	}
	return *o.PollingPeriod
}

// GetPollingPeriodOk returns a tuple with the PollingPeriod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeConfig) GetPollingPeriodOk() (*string, bool) {
	if o == nil || isNil(o.PollingPeriod) {
		return nil, false
	}
	return o.PollingPeriod, true
}

// HasPollingPeriod returns a boolean if a field has been set.
func (o *StripeConfig) HasPollingPeriod() bool {
	if o != nil && !isNil(o.PollingPeriod) {
		return true
	}

	return false
}

// SetPollingPeriod gets a reference to the given string and assigns it to the PollingPeriod field.
func (o *StripeConfig) SetPollingPeriod(v string) {
	o.PollingPeriod = &v
}

// GetApiKey returns the ApiKey field value
func (o *StripeConfig) GetApiKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value
// and a boolean to check if the value has been set.
func (o *StripeConfig) GetApiKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApiKey, true
}

// SetApiKey sets field value
func (o *StripeConfig) SetApiKey(v string) {
	o.ApiKey = v
}

// GetPageSize returns the PageSize field value if set, zero value otherwise.
func (o *StripeConfig) GetPageSize() float32 {
	if o == nil || isNil(o.PageSize) {
		var ret float32
		return ret
	}
	return *o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StripeConfig) GetPageSizeOk() (*float32, bool) {
	if o == nil || isNil(o.PageSize) {
		return nil, false
	}
	return o.PageSize, true
}

// HasPageSize returns a boolean if a field has been set.
func (o *StripeConfig) HasPageSize() bool {
	if o != nil && !isNil(o.PageSize) {
		return true
	}

	return false
}

// SetPageSize gets a reference to the given float32 and assigns it to the PageSize field.
func (o *StripeConfig) SetPageSize(v float32) {
	o.PageSize = &v
}

func (o StripeConfig) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StripeConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.PollingPeriod) {
		toSerialize["pollingPeriod"] = o.PollingPeriod
	}
	toSerialize["apiKey"] = o.ApiKey
	if !isNil(o.PageSize) {
		toSerialize["pageSize"] = o.PageSize
	}
	return toSerialize, nil
}

type NullableStripeConfig struct {
	value *StripeConfig
	isSet bool
}

func (v NullableStripeConfig) Get() *StripeConfig {
	return v.value
}

func (v *NullableStripeConfig) Set(val *StripeConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableStripeConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableStripeConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStripeConfig(val *StripeConfig) *NullableStripeConfig {
	return &NullableStripeConfig{value: val, isSet: true}
}

func (v NullableStripeConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStripeConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
