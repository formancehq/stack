/*
Auth API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: AUTH_VERSION
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package authclient

import (
	"encoding/json"
)

// checks if the Secret type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Secret{}

// Secret struct for Secret
type Secret struct {
	Name string `json:"name"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Id string `json:"id"`
	LastDigits string `json:"lastDigits"`
	Clear string `json:"clear"`
}

// NewSecret instantiates a new Secret object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecret(name string, id string, lastDigits string, clear string) *Secret {
	this := Secret{}
	this.Name = name
	this.Id = id
	this.LastDigits = lastDigits
	this.Clear = clear
	return &this
}

// NewSecretWithDefaults instantiates a new Secret object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecretWithDefaults() *Secret {
	this := Secret{}
	return &this
}

// GetName returns the Name field value
func (o *Secret) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Secret) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Secret) SetName(v string) {
	o.Name = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *Secret) GetMetadata() map[string]interface{} {
	if o == nil || isNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Secret) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || isNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *Secret) HasMetadata() bool {
	if o != nil && !isNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *Secret) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetId returns the Id field value
func (o *Secret) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Secret) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Secret) SetId(v string) {
	o.Id = v
}

// GetLastDigits returns the LastDigits field value
func (o *Secret) GetLastDigits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LastDigits
}

// GetLastDigitsOk returns a tuple with the LastDigits field value
// and a boolean to check if the value has been set.
func (o *Secret) GetLastDigitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastDigits, true
}

// SetLastDigits sets field value
func (o *Secret) SetLastDigits(v string) {
	o.LastDigits = v
}

// GetClear returns the Clear field value
func (o *Secret) GetClear() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Clear
}

// GetClearOk returns a tuple with the Clear field value
// and a boolean to check if the value has been set.
func (o *Secret) GetClearOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Clear, true
}

// SetClear sets field value
func (o *Secret) SetClear(v string) {
	o.Clear = v
}

func (o Secret) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Secret) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !isNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	toSerialize["id"] = o.Id
	toSerialize["lastDigits"] = o.LastDigits
	toSerialize["clear"] = o.Clear
	return toSerialize, nil
}

type NullableSecret struct {
	value *Secret
	isSet bool
}

func (v NullableSecret) Get() *Secret {
	return v.value
}

func (v *NullableSecret) Set(val *Secret) {
	v.value = val
	v.isSet = true
}

func (v NullableSecret) IsSet() bool {
	return v.isSet
}

func (v *NullableSecret) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecret(val *Secret) *NullableSecret {
	return &NullableSecret{value: val, isSet: true}
}

func (v NullableSecret) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecret) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


