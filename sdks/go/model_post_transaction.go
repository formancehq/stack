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
	"time"
)

// checks if the PostTransaction type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PostTransaction{}

// PostTransaction struct for PostTransaction
type PostTransaction struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Postings []Posting `json:"postings,omitempty"`
	Script *PostTransactionScript `json:"script,omitempty"`
	Reference *string `json:"reference,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewPostTransaction instantiates a new PostTransaction object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPostTransaction() *PostTransaction {
	this := PostTransaction{}
	return &this
}

// NewPostTransactionWithDefaults instantiates a new PostTransaction object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPostTransactionWithDefaults() *PostTransaction {
	this := PostTransaction{}
	return &this
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise.
func (o *PostTransaction) GetTimestamp() time.Time {
	if o == nil || IsNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostTransaction) GetTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}
	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *PostTransaction) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *PostTransaction) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}

// GetPostings returns the Postings field value if set, zero value otherwise.
func (o *PostTransaction) GetPostings() []Posting {
	if o == nil || IsNil(o.Postings) {
		var ret []Posting
		return ret
	}
	return o.Postings
}

// GetPostingsOk returns a tuple with the Postings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostTransaction) GetPostingsOk() ([]Posting, bool) {
	if o == nil || IsNil(o.Postings) {
		return nil, false
	}
	return o.Postings, true
}

// HasPostings returns a boolean if a field has been set.
func (o *PostTransaction) HasPostings() bool {
	if o != nil && !IsNil(o.Postings) {
		return true
	}

	return false
}

// SetPostings gets a reference to the given []Posting and assigns it to the Postings field.
func (o *PostTransaction) SetPostings(v []Posting) {
	o.Postings = v
}

// GetScript returns the Script field value if set, zero value otherwise.
func (o *PostTransaction) GetScript() PostTransactionScript {
	if o == nil || IsNil(o.Script) {
		var ret PostTransactionScript
		return ret
	}
	return *o.Script
}

// GetScriptOk returns a tuple with the Script field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostTransaction) GetScriptOk() (*PostTransactionScript, bool) {
	if o == nil || IsNil(o.Script) {
		return nil, false
	}
	return o.Script, true
}

// HasScript returns a boolean if a field has been set.
func (o *PostTransaction) HasScript() bool {
	if o != nil && !IsNil(o.Script) {
		return true
	}

	return false
}

// SetScript gets a reference to the given PostTransactionScript and assigns it to the Script field.
func (o *PostTransaction) SetScript(v PostTransactionScript) {
	o.Script = &v
}

// GetReference returns the Reference field value if set, zero value otherwise.
func (o *PostTransaction) GetReference() string {
	if o == nil || IsNil(o.Reference) {
		var ret string
		return ret
	}
	return *o.Reference
}

// GetReferenceOk returns a tuple with the Reference field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostTransaction) GetReferenceOk() (*string, bool) {
	if o == nil || IsNil(o.Reference) {
		return nil, false
	}
	return o.Reference, true
}

// HasReference returns a boolean if a field has been set.
func (o *PostTransaction) HasReference() bool {
	if o != nil && !IsNil(o.Reference) {
		return true
	}

	return false
}

// SetReference gets a reference to the given string and assigns it to the Reference field.
func (o *PostTransaction) SetReference(v string) {
	o.Reference = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PostTransaction) GetMetadata() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PostTransaction) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *PostTransaction) HasMetadata() bool {
	if o != nil && IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *PostTransaction) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

func (o PostTransaction) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PostTransaction) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Timestamp) {
		toSerialize["timestamp"] = o.Timestamp
	}
	if !IsNil(o.Postings) {
		toSerialize["postings"] = o.Postings
	}
	if !IsNil(o.Script) {
		toSerialize["script"] = o.Script
	}
	if !IsNil(o.Reference) {
		toSerialize["reference"] = o.Reference
	}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullablePostTransaction struct {
	value *PostTransaction
	isSet bool
}

func (v NullablePostTransaction) Get() *PostTransaction {
	return v.value
}

func (v *NullablePostTransaction) Set(val *PostTransaction) {
	v.value = val
	v.isSet = true
}

func (v NullablePostTransaction) IsSet() bool {
	return v.isSet
}

func (v *NullablePostTransaction) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePostTransaction(val *PostTransaction) *NullablePostTransaction {
	return &NullablePostTransaction{value: val, isSet: true}
}

func (v NullablePostTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePostTransaction) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


