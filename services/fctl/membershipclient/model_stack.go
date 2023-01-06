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

// checks if the Stack type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Stack{}

// Stack struct for Stack
type Stack struct {
	// Stack name
	Name string `json:"name"`
	Tags map[string]string `json:"tags"`
	Production bool `json:"production"`
	Metadata map[string]string `json:"metadata"`
	// Stack ID
	Id string `json:"id"`
	// Organization ID
	OrganizationId string `json:"organizationId"`
	// Base stack uri
	Uri string `json:"uri"`
	BoundRegion *Region `json:"boundRegion,omitempty"`
}

// NewStack instantiates a new Stack object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStack(name string, tags map[string]string, production bool, metadata map[string]string, id string, organizationId string, uri string) *Stack {
	this := Stack{}
	this.Name = name
	this.Tags = tags
	this.Production = production
	this.Metadata = metadata
	this.Id = id
	this.OrganizationId = organizationId
	this.Uri = uri
	return &this
}

// NewStackWithDefaults instantiates a new Stack object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStackWithDefaults() *Stack {
	this := Stack{}
	return &this
}

// GetName returns the Name field value
func (o *Stack) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Stack) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Stack) SetName(v string) {
	o.Name = v
}

// GetTags returns the Tags field value
func (o *Stack) GetTags() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value
// and a boolean to check if the value has been set.
func (o *Stack) GetTagsOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Tags, true
}

// SetTags sets field value
func (o *Stack) SetTags(v map[string]string) {
	o.Tags = v
}

// GetProduction returns the Production field value
func (o *Stack) GetProduction() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Production
}

// GetProductionOk returns a tuple with the Production field value
// and a boolean to check if the value has been set.
func (o *Stack) GetProductionOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Production, true
}

// SetProduction sets field value
func (o *Stack) SetProduction(v bool) {
	o.Production = v
}

// GetMetadata returns the Metadata field value
func (o *Stack) GetMetadata() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *Stack) GetMetadataOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *Stack) SetMetadata(v map[string]string) {
	o.Metadata = v
}

// GetId returns the Id field value
func (o *Stack) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Stack) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Stack) SetId(v string) {
	o.Id = v
}

// GetOrganizationId returns the OrganizationId field value
func (o *Stack) GetOrganizationId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value
// and a boolean to check if the value has been set.
func (o *Stack) GetOrganizationIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrganizationId, true
}

// SetOrganizationId sets field value
func (o *Stack) SetOrganizationId(v string) {
	o.OrganizationId = v
}

// GetUri returns the Uri field value
func (o *Stack) GetUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uri
}

// GetUriOk returns a tuple with the Uri field value
// and a boolean to check if the value has been set.
func (o *Stack) GetUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uri, true
}

// SetUri sets field value
func (o *Stack) SetUri(v string) {
	o.Uri = v
}

// GetBoundRegion returns the BoundRegion field value if set, zero value otherwise.
func (o *Stack) GetBoundRegion() Region {
	if o == nil || isNil(o.BoundRegion) {
		var ret Region
		return ret
	}
	return *o.BoundRegion
}

// GetBoundRegionOk returns a tuple with the BoundRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Stack) GetBoundRegionOk() (*Region, bool) {
	if o == nil || isNil(o.BoundRegion) {
		return nil, false
	}
	return o.BoundRegion, true
}

// HasBoundRegion returns a boolean if a field has been set.
func (o *Stack) HasBoundRegion() bool {
	if o != nil && !isNil(o.BoundRegion) {
		return true
	}

	return false
}

// SetBoundRegion gets a reference to the given Region and assigns it to the BoundRegion field.
func (o *Stack) SetBoundRegion(v Region) {
	o.BoundRegion = &v
}

func (o Stack) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Stack) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["tags"] = o.Tags
	toSerialize["production"] = o.Production
	toSerialize["metadata"] = o.Metadata
	toSerialize["id"] = o.Id
	toSerialize["organizationId"] = o.OrganizationId
	toSerialize["uri"] = o.Uri
	if !isNil(o.BoundRegion) {
		toSerialize["boundRegion"] = o.BoundRegion
	}
	return toSerialize, nil
}

type NullableStack struct {
	value *Stack
	isSet bool
}

func (v NullableStack) Get() *Stack {
	return v.value
}

func (v *NullableStack) Set(val *Stack) {
	v.value = val
	v.isSet = true
}

func (v NullableStack) IsSet() bool {
	return v.isSet
}

func (v *NullableStack) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStack(val *Stack) *NullableStack {
	return &NullableStack{value: val, isSet: true}
}

func (v NullableStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStack) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
