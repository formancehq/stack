/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230301
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the ScopeOptions type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ScopeOptions{}

// ScopeOptions struct for ScopeOptions
type ScopeOptions struct {
	Label string `json:"label"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewScopeOptions instantiates a new ScopeOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScopeOptions(label string) *ScopeOptions {
	this := ScopeOptions{}
	this.Label = label
	return &this
}

// NewScopeOptionsWithDefaults instantiates a new ScopeOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScopeOptionsWithDefaults() *ScopeOptions {
	this := ScopeOptions{}
	return &this
}

// GetLabel returns the Label field value
func (o *ScopeOptions) GetLabel() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Label
}

// GetLabelOk returns a tuple with the Label field value
// and a boolean to check if the value has been set.
func (o *ScopeOptions) GetLabelOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Label, true
}

// SetLabel sets field value
func (o *ScopeOptions) SetLabel(v string) {
	o.Label = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *ScopeOptions) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeOptions) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *ScopeOptions) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *ScopeOptions) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

func (o ScopeOptions) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ScopeOptions) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["label"] = o.Label
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullableScopeOptions struct {
	value *ScopeOptions
	isSet bool
}

func (v NullableScopeOptions) Get() *ScopeOptions {
	return v.value
}

func (v *NullableScopeOptions) Set(val *ScopeOptions) {
	v.value = val
	v.isSet = true
}

func (v NullableScopeOptions) IsSet() bool {
	return v.isSet
}

func (v *NullableScopeOptions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScopeOptions(val *ScopeOptions) *NullableScopeOptions {
	return &NullableScopeOptions{value: val, isSet: true}
}

func (v NullableScopeOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScopeOptions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


