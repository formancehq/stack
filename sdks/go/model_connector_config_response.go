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

// ConnectorConfigResponse struct for ConnectorConfigResponse
type ConnectorConfigResponse struct {
	Data *ConnectorConfig `json:"data,omitempty"`
}

// NewConnectorConfigResponse instantiates a new ConnectorConfigResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectorConfigResponse() *ConnectorConfigResponse {
	this := ConnectorConfigResponse{}
	return &this
}

// NewConnectorConfigResponseWithDefaults instantiates a new ConnectorConfigResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectorConfigResponseWithDefaults() *ConnectorConfigResponse {
	this := ConnectorConfigResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ConnectorConfigResponse) GetData() ConnectorConfig {
	if o == nil || isNil(o.Data) {
		var ret ConnectorConfig
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectorConfigResponse) GetDataOk() (*ConnectorConfig, bool) {
	if o == nil || isNil(o.Data) {
    return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ConnectorConfigResponse) HasData() bool {
	if o != nil && !isNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given ConnectorConfig and assigns it to the Data field.
func (o *ConnectorConfigResponse) SetData(v ConnectorConfig) {
	o.Data = &v
}

func (o ConnectorConfigResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableConnectorConfigResponse struct {
	value *ConnectorConfigResponse
	isSet bool
}

func (v NullableConnectorConfigResponse) Get() *ConnectorConfigResponse {
	return v.value
}

func (v *NullableConnectorConfigResponse) Set(val *ConnectorConfigResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectorConfigResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectorConfigResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectorConfigResponse(val *ConnectorConfigResponse) *NullableConnectorConfigResponse {
	return &NullableConnectorConfigResponse{value: val, isSet: true}
}

func (v NullableConnectorConfigResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectorConfigResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


