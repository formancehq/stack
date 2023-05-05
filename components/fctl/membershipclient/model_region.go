/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
	"time"
)

// checks if the Region type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Region{}

// Region struct for Region
type Region struct {
	Id        string     `json:"id"`
	BaseUrl   string     `json:"baseUrl"`
	CreatedAt string     `json:"createdAt"`
	Active    bool       `json:"active"`
	LastPing  *time.Time `json:"lastPing,omitempty"`
	Name      string     `json:"name"`
}

// NewRegion instantiates a new Region object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegion(id string, baseUrl string, createdAt string, active bool, name string) *Region {
	this := Region{}
	this.Id = id
	this.BaseUrl = baseUrl
	this.CreatedAt = createdAt
	this.Active = active
	this.Name = name
	return &this
}

// NewRegionWithDefaults instantiates a new Region object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegionWithDefaults() *Region {
	this := Region{}
	return &this
}

// GetId returns the Id field value
func (o *Region) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Region) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Region) SetId(v string) {
	o.Id = v
}

// GetBaseUrl returns the BaseUrl field value
func (o *Region) GetBaseUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BaseUrl
}

// GetBaseUrlOk returns a tuple with the BaseUrl field value
// and a boolean to check if the value has been set.
func (o *Region) GetBaseUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BaseUrl, true
}

// SetBaseUrl sets field value
func (o *Region) SetBaseUrl(v string) {
	o.BaseUrl = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Region) GetCreatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Region) GetCreatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Region) SetCreatedAt(v string) {
	o.CreatedAt = v
}

// GetActive returns the Active field value
func (o *Region) GetActive() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Active
}

// GetActiveOk returns a tuple with the Active field value
// and a boolean to check if the value has been set.
func (o *Region) GetActiveOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Active, true
}

// SetActive sets field value
func (o *Region) SetActive(v bool) {
	o.Active = v
}

// GetLastPing returns the LastPing field value if set, zero value otherwise.
func (o *Region) GetLastPing() time.Time {
	if o == nil || IsNil(o.LastPing) {
		var ret time.Time
		return ret
	}
	return *o.LastPing
}

// GetLastPingOk returns a tuple with the LastPing field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Region) GetLastPingOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastPing) {
		return nil, false
	}
	return o.LastPing, true
}

// HasLastPing returns a boolean if a field has been set.
func (o *Region) HasLastPing() bool {
	if o != nil && !IsNil(o.LastPing) {
		return true
	}

	return false
}

// SetLastPing gets a reference to the given time.Time and assigns it to the LastPing field.
func (o *Region) SetLastPing(v time.Time) {
	o.LastPing = &v
}

// GetName returns the Name field value
func (o *Region) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Region) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Region) SetName(v string) {
	o.Name = v
}

func (o Region) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Region) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["baseUrl"] = o.BaseUrl
	toSerialize["createdAt"] = o.CreatedAt
	toSerialize["active"] = o.Active
	if !IsNil(o.LastPing) {
		toSerialize["lastPing"] = o.LastPing
	}
	toSerialize["name"] = o.Name
	return toSerialize, nil
}

type NullableRegion struct {
	value *Region
	isSet bool
}

func (v NullableRegion) Get() *Region {
	return v.value
}

func (v *NullableRegion) Set(val *Region) {
	v.value = val
	v.isSet = true
}

func (v NullableRegion) IsSet() bool {
	return v.isSet
}

func (v *NullableRegion) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegion(val *Region) *NullableRegion {
	return &NullableRegion{value: val, isSet: true}
}

func (v NullableRegion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegion) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
