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

// ListClientsResponse struct for ListClientsResponse
type ListClientsResponse struct {
	Data []Client `json:"data,omitempty"`
}

// NewListClientsResponse instantiates a new ListClientsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListClientsResponse() *ListClientsResponse {
	this := ListClientsResponse{}
	return &this
}

// NewListClientsResponseWithDefaults instantiates a new ListClientsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListClientsResponseWithDefaults() *ListClientsResponse {
	this := ListClientsResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ListClientsResponse) GetData() []Client {
	if o == nil || isNil(o.Data) {
		var ret []Client
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListClientsResponse) GetDataOk() ([]Client, bool) {
	if o == nil || isNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ListClientsResponse) HasData() bool {
	if o != nil && !isNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given []Client and assigns it to the Data field.
func (o *ListClientsResponse) SetData(v []Client) {
	o.Data = v
}

func (o ListClientsResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableListClientsResponse struct {
	value *ListClientsResponse
	isSet bool
}

func (v NullableListClientsResponse) Get() *ListClientsResponse {
	return v.value
}

func (v *NullableListClientsResponse) Set(val *ListClientsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableListClientsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableListClientsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListClientsResponse(val *ListClientsResponse) *NullableListClientsResponse {
	return &NullableListClientsResponse{value: val, isSet: true}
}

func (v NullableListClientsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListClientsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
