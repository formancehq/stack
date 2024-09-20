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

// checks if the StackAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StackAllOf{}

// StackAllOf struct for StackAllOf
type StackAllOf struct {
	// Stack ID
	Id string `json:"id"`
	// Organization ID
	OrganizationId string `json:"organizationId"`
	// Base stack uri
	Uri string `json:"uri"`
	// The region where the stack is installed
	RegionID string `json:"regionID"`
	StargateEnabled bool `json:"stargateEnabled"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	DisabledAt *time.Time `json:"disabledAt,omitempty"`
	AuditEnabled *bool `json:"auditEnabled,omitempty"`
	Synchronised bool `json:"synchronised"`
}

// NewStackAllOf instantiates a new StackAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStackAllOf(id string, organizationId string, uri string, regionID string, stargateEnabled bool, synchronised bool) *StackAllOf {
	this := StackAllOf{}
	this.Id = id
	this.OrganizationId = organizationId
	this.Uri = uri
	this.RegionID = regionID
	this.StargateEnabled = stargateEnabled
	this.Synchronised = synchronised
	return &this
}

// NewStackAllOfWithDefaults instantiates a new StackAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStackAllOfWithDefaults() *StackAllOf {
	this := StackAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *StackAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *StackAllOf) SetId(v string) {
	o.Id = v
}

// GetOrganizationId returns the OrganizationId field value
func (o *StackAllOf) GetOrganizationId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetOrganizationIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrganizationId, true
}

// SetOrganizationId sets field value
func (o *StackAllOf) SetOrganizationId(v string) {
	o.OrganizationId = v
}

// GetUri returns the Uri field value
func (o *StackAllOf) GetUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uri
}

// GetUriOk returns a tuple with the Uri field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uri, true
}

// SetUri sets field value
func (o *StackAllOf) SetUri(v string) {
	o.Uri = v
}

// GetRegionID returns the RegionID field value
func (o *StackAllOf) GetRegionID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionID
}

// GetRegionIDOk returns a tuple with the RegionID field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetRegionIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionID, true
}

// SetRegionID sets field value
func (o *StackAllOf) SetRegionID(v string) {
	o.RegionID = v
}

// GetStargateEnabled returns the StargateEnabled field value
func (o *StackAllOf) GetStargateEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.StargateEnabled
}

// GetStargateEnabledOk returns a tuple with the StargateEnabled field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetStargateEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StargateEnabled, true
}

// SetStargateEnabled sets field value
func (o *StackAllOf) SetStargateEnabled(v bool) {
	o.StargateEnabled = v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *StackAllOf) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *StackAllOf) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *StackAllOf) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetDeletedAt returns the DeletedAt field value if set, zero value otherwise.
func (o *StackAllOf) GetDeletedAt() time.Time {
	if o == nil || IsNil(o.DeletedAt) {
		var ret time.Time
		return ret
	}
	return *o.DeletedAt
}

// GetDeletedAtOk returns a tuple with the DeletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetDeletedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DeletedAt) {
		return nil, false
	}
	return o.DeletedAt, true
}

// HasDeletedAt returns a boolean if a field has been set.
func (o *StackAllOf) HasDeletedAt() bool {
	if o != nil && !IsNil(o.DeletedAt) {
		return true
	}

	return false
}

// SetDeletedAt gets a reference to the given time.Time and assigns it to the DeletedAt field.
func (o *StackAllOf) SetDeletedAt(v time.Time) {
	o.DeletedAt = &v
}

// GetDisabledAt returns the DisabledAt field value if set, zero value otherwise.
func (o *StackAllOf) GetDisabledAt() time.Time {
	if o == nil || IsNil(o.DisabledAt) {
		var ret time.Time
		return ret
	}
	return *o.DisabledAt
}

// GetDisabledAtOk returns a tuple with the DisabledAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetDisabledAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DisabledAt) {
		return nil, false
	}
	return o.DisabledAt, true
}

// HasDisabledAt returns a boolean if a field has been set.
func (o *StackAllOf) HasDisabledAt() bool {
	if o != nil && !IsNil(o.DisabledAt) {
		return true
	}

	return false
}

// SetDisabledAt gets a reference to the given time.Time and assigns it to the DisabledAt field.
func (o *StackAllOf) SetDisabledAt(v time.Time) {
	o.DisabledAt = &v
}

// GetAuditEnabled returns the AuditEnabled field value if set, zero value otherwise.
func (o *StackAllOf) GetAuditEnabled() bool {
	if o == nil || IsNil(o.AuditEnabled) {
		var ret bool
		return ret
	}
	return *o.AuditEnabled
}

// GetAuditEnabledOk returns a tuple with the AuditEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetAuditEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AuditEnabled) {
		return nil, false
	}
	return o.AuditEnabled, true
}

// HasAuditEnabled returns a boolean if a field has been set.
func (o *StackAllOf) HasAuditEnabled() bool {
	if o != nil && !IsNil(o.AuditEnabled) {
		return true
	}

	return false
}

// SetAuditEnabled gets a reference to the given bool and assigns it to the AuditEnabled field.
func (o *StackAllOf) SetAuditEnabled(v bool) {
	o.AuditEnabled = &v
}

// GetSynchronised returns the Synchronised field value
func (o *StackAllOf) GetSynchronised() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Synchronised
}

// GetSynchronisedOk returns a tuple with the Synchronised field value
// and a boolean to check if the value has been set.
func (o *StackAllOf) GetSynchronisedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Synchronised, true
}

// SetSynchronised sets field value
func (o *StackAllOf) SetSynchronised(v bool) {
	o.Synchronised = v
}

func (o StackAllOf) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StackAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["organizationId"] = o.OrganizationId
	toSerialize["uri"] = o.Uri
	toSerialize["regionID"] = o.RegionID
	toSerialize["stargateEnabled"] = o.StargateEnabled
	if !IsNil(o.CreatedAt) {
		toSerialize["createdAt"] = o.CreatedAt
	}
	if !IsNil(o.DeletedAt) {
		toSerialize["deletedAt"] = o.DeletedAt
	}
	if !IsNil(o.DisabledAt) {
		toSerialize["disabledAt"] = o.DisabledAt
	}
	if !IsNil(o.AuditEnabled) {
		toSerialize["auditEnabled"] = o.AuditEnabled
	}
	toSerialize["synchronised"] = o.Synchronised
	return toSerialize, nil
}

type NullableStackAllOf struct {
	value *StackAllOf
	isSet bool
}

func (v NullableStackAllOf) Get() *StackAllOf {
	return v.value
}

func (v *NullableStackAllOf) Set(val *StackAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableStackAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableStackAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStackAllOf(val *StackAllOf) *NullableStackAllOf {
	return &NullableStackAllOf{value: val, isSet: true}
}

func (v NullableStackAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStackAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


