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

// checks if the GetRegionVersionsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetRegionVersionsResponse{}

// GetRegionVersionsResponse struct for GetRegionVersionsResponse
type GetRegionVersionsResponse struct {
	Data []Version `json:"data"`
}

// NewGetRegionVersionsResponse instantiates a new GetRegionVersionsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetRegionVersionsResponse(data []Version) *GetRegionVersionsResponse {
	this := GetRegionVersionsResponse{}
	this.Data = data
	return &this
}

// NewGetRegionVersionsResponseWithDefaults instantiates a new GetRegionVersionsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetRegionVersionsResponseWithDefaults() *GetRegionVersionsResponse {
	this := GetRegionVersionsResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetRegionVersionsResponse) GetData() []Version {
	if o == nil {
		var ret []Version
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetRegionVersionsResponse) GetDataOk() ([]Version, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *GetRegionVersionsResponse) SetData(v []Version) {
	o.Data = v
}

func (o GetRegionVersionsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetRegionVersionsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableGetRegionVersionsResponse struct {
	value *GetRegionVersionsResponse
	isSet bool
}

func (v NullableGetRegionVersionsResponse) Get() *GetRegionVersionsResponse {
	return v.value
}

func (v *NullableGetRegionVersionsResponse) Set(val *GetRegionVersionsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetRegionVersionsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetRegionVersionsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetRegionVersionsResponse(val *GetRegionVersionsResponse) *NullableGetRegionVersionsResponse {
	return &NullableGetRegionVersionsResponse{value: val, isSet: true}
}

func (v NullableGetRegionVersionsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetRegionVersionsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


