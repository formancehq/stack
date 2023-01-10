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

// ConnectorsResponseDataInner struct for ConnectorsResponseDataInner
type ConnectorsResponseDataInner struct {
	Provider *Connector `json:"provider,omitempty"`
	Enabled *bool `json:"enabled,omitempty"`
	Disabled *bool `json:"disabled,omitempty"`
}

// NewConnectorsResponseDataInner instantiates a new ConnectorsResponseDataInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectorsResponseDataInner() *ConnectorsResponseDataInner {
	this := ConnectorsResponseDataInner{}
	return &this
}

// NewConnectorsResponseDataInnerWithDefaults instantiates a new ConnectorsResponseDataInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectorsResponseDataInnerWithDefaults() *ConnectorsResponseDataInner {
	this := ConnectorsResponseDataInner{}
	return &this
}

// GetProvider returns the Provider field value if set, zero value otherwise.
func (o *ConnectorsResponseDataInner) GetProvider() Connector {
	if o == nil || isNil(o.Provider) {
		var ret Connector
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectorsResponseDataInner) GetProviderOk() (*Connector, bool) {
	if o == nil || isNil(o.Provider) {
    return nil, false
	}
	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *ConnectorsResponseDataInner) HasProvider() bool {
	if o != nil && !isNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given Connector and assigns it to the Provider field.
func (o *ConnectorsResponseDataInner) SetProvider(v Connector) {
	o.Provider = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *ConnectorsResponseDataInner) GetEnabled() bool {
	if o == nil || isNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectorsResponseDataInner) GetEnabledOk() (*bool, bool) {
	if o == nil || isNil(o.Enabled) {
    return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *ConnectorsResponseDataInner) HasEnabled() bool {
	if o != nil && !isNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *ConnectorsResponseDataInner) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetDisabled returns the Disabled field value if set, zero value otherwise.
func (o *ConnectorsResponseDataInner) GetDisabled() bool {
	if o == nil || isNil(o.Disabled) {
		var ret bool
		return ret
	}
	return *o.Disabled
}

// GetDisabledOk returns a tuple with the Disabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectorsResponseDataInner) GetDisabledOk() (*bool, bool) {
	if o == nil || isNil(o.Disabled) {
    return nil, false
	}
	return o.Disabled, true
}

// HasDisabled returns a boolean if a field has been set.
func (o *ConnectorsResponseDataInner) HasDisabled() bool {
	if o != nil && !isNil(o.Disabled) {
		return true
	}

	return false
}

// SetDisabled gets a reference to the given bool and assigns it to the Disabled field.
func (o *ConnectorsResponseDataInner) SetDisabled(v bool) {
	o.Disabled = &v
}

func (o ConnectorsResponseDataInner) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Provider) {
		toSerialize["provider"] = o.Provider
	}
	if !isNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	if !isNil(o.Disabled) {
		toSerialize["disabled"] = o.Disabled
	}
	return json.Marshal(toSerialize)
}

type NullableConnectorsResponseDataInner struct {
	value *ConnectorsResponseDataInner
	isSet bool
}

func (v NullableConnectorsResponseDataInner) Get() *ConnectorsResponseDataInner {
	return v.value
}

func (v *NullableConnectorsResponseDataInner) Set(val *ConnectorsResponseDataInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectorsResponseDataInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectorsResponseDataInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectorsResponseDataInner(val *ConnectorsResponseDataInner) *NullableConnectorsResponseDataInner {
	return &NullableConnectorsResponseDataInner{value: val, isSet: true}
}

func (v NullableConnectorsResponseDataInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectorsResponseDataInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


