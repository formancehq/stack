/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
)

// checks if the CreateInvitationResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateInvitationResponse{}

// CreateInvitationResponse struct for CreateInvitationResponse
type CreateInvitationResponse struct {
	Data *Invitation `json:"data,omitempty"`
}

// NewCreateInvitationResponse instantiates a new CreateInvitationResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateInvitationResponse() *CreateInvitationResponse {
	this := CreateInvitationResponse{}
	return &this
}

// NewCreateInvitationResponseWithDefaults instantiates a new CreateInvitationResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateInvitationResponseWithDefaults() *CreateInvitationResponse {
	this := CreateInvitationResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *CreateInvitationResponse) GetData() Invitation {
	if o == nil || isNil(o.Data) {
		var ret Invitation
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateInvitationResponse) GetDataOk() (*Invitation, bool) {
	if o == nil || isNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *CreateInvitationResponse) HasData() bool {
	if o != nil && !isNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given Invitation and assigns it to the Data field.
func (o *CreateInvitationResponse) SetData(v Invitation) {
	o.Data = &v
}

func (o CreateInvitationResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateInvitationResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableCreateInvitationResponse struct {
	value *CreateInvitationResponse
	isSet bool
}

func (v NullableCreateInvitationResponse) Get() *CreateInvitationResponse {
	return v.value
}

func (v *NullableCreateInvitationResponse) Set(val *CreateInvitationResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateInvitationResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateInvitationResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateInvitationResponse(val *CreateInvitationResponse) *NullableCreateInvitationResponse {
	return &NullableCreateInvitationResponse{value: val, isSet: true}
}

func (v NullableCreateInvitationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateInvitationResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
