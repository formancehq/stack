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

// checks if the ExpandedTransaction type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExpandedTransaction{}

// ExpandedTransaction struct for ExpandedTransaction
type ExpandedTransaction struct {
	Timestamp time.Time `json:"timestamp"`
	Postings []Posting `json:"postings"`
	Reference *string `json:"reference,omitempty"`
	Metadata map[string]string `json:"metadata"`
	Txid int64 `json:"txid"`
	PreCommitVolumes map[string]map[string]Volume `json:"preCommitVolumes"`
	PostCommitVolumes map[string]map[string]Volume `json:"postCommitVolumes"`
}

// NewExpandedTransaction instantiates a new ExpandedTransaction object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExpandedTransaction(timestamp time.Time, postings []Posting, metadata map[string]string, txid int64, preCommitVolumes map[string]map[string]Volume, postCommitVolumes map[string]map[string]Volume) *ExpandedTransaction {
	this := ExpandedTransaction{}
	this.Timestamp = timestamp
	this.Postings = postings
	this.Metadata = metadata
	this.Txid = txid
	this.PreCommitVolumes = preCommitVolumes
	this.PostCommitVolumes = postCommitVolumes
	return &this
}

// NewExpandedTransactionWithDefaults instantiates a new ExpandedTransaction object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExpandedTransactionWithDefaults() *ExpandedTransaction {
	this := ExpandedTransaction{}
	return &this
}

// GetTimestamp returns the Timestamp field value
func (o *ExpandedTransaction) GetTimestamp() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetTimestampOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *ExpandedTransaction) SetTimestamp(v time.Time) {
	o.Timestamp = v
}

// GetPostings returns the Postings field value
func (o *ExpandedTransaction) GetPostings() []Posting {
	if o == nil {
		var ret []Posting
		return ret
	}

	return o.Postings
}

// GetPostingsOk returns a tuple with the Postings field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetPostingsOk() ([]Posting, bool) {
	if o == nil {
		return nil, false
	}
	return o.Postings, true
}

// SetPostings sets field value
func (o *ExpandedTransaction) SetPostings(v []Posting) {
	o.Postings = v
}

// GetReference returns the Reference field value if set, zero value otherwise.
func (o *ExpandedTransaction) GetReference() string {
	if o == nil || IsNil(o.Reference) {
		var ret string
		return ret
	}
	return *o.Reference
}

// GetReferenceOk returns a tuple with the Reference field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetReferenceOk() (*string, bool) {
	if o == nil || IsNil(o.Reference) {
		return nil, false
	}
	return o.Reference, true
}

// HasReference returns a boolean if a field has been set.
func (o *ExpandedTransaction) HasReference() bool {
	if o != nil && !IsNil(o.Reference) {
		return true
	}

	return false
}

// SetReference gets a reference to the given string and assigns it to the Reference field.
func (o *ExpandedTransaction) SetReference(v string) {
	o.Reference = &v
}

// GetMetadata returns the Metadata field value
func (o *ExpandedTransaction) GetMetadata() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetMetadataOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *ExpandedTransaction) SetMetadata(v map[string]string) {
	o.Metadata = v
}

// GetTxid returns the Txid field value
func (o *ExpandedTransaction) GetTxid() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Txid
}

// GetTxidOk returns a tuple with the Txid field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetTxidOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Txid, true
}

// SetTxid sets field value
func (o *ExpandedTransaction) SetTxid(v int64) {
	o.Txid = v
}

// GetPreCommitVolumes returns the PreCommitVolumes field value
func (o *ExpandedTransaction) GetPreCommitVolumes() map[string]map[string]Volume {
	if o == nil {
		var ret map[string]map[string]Volume
		return ret
	}

	return o.PreCommitVolumes
}

// GetPreCommitVolumesOk returns a tuple with the PreCommitVolumes field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetPreCommitVolumesOk() (*map[string]map[string]Volume, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PreCommitVolumes, true
}

// SetPreCommitVolumes sets field value
func (o *ExpandedTransaction) SetPreCommitVolumes(v map[string]map[string]Volume) {
	o.PreCommitVolumes = v
}

// GetPostCommitVolumes returns the PostCommitVolumes field value
func (o *ExpandedTransaction) GetPostCommitVolumes() map[string]map[string]Volume {
	if o == nil {
		var ret map[string]map[string]Volume
		return ret
	}

	return o.PostCommitVolumes
}

// GetPostCommitVolumesOk returns a tuple with the PostCommitVolumes field value
// and a boolean to check if the value has been set.
func (o *ExpandedTransaction) GetPostCommitVolumesOk() (*map[string]map[string]Volume, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PostCommitVolumes, true
}

// SetPostCommitVolumes sets field value
func (o *ExpandedTransaction) SetPostCommitVolumes(v map[string]map[string]Volume) {
	o.PostCommitVolumes = v
}

func (o ExpandedTransaction) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExpandedTransaction) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["timestamp"] = o.Timestamp
	toSerialize["postings"] = o.Postings
	if !IsNil(o.Reference) {
		toSerialize["reference"] = o.Reference
	}
	toSerialize["metadata"] = o.Metadata
	toSerialize["txid"] = o.Txid
	toSerialize["preCommitVolumes"] = o.PreCommitVolumes
	toSerialize["postCommitVolumes"] = o.PostCommitVolumes
	return toSerialize, nil
}

type NullableExpandedTransaction struct {
	value *ExpandedTransaction
	isSet bool
}

func (v NullableExpandedTransaction) Get() *ExpandedTransaction {
	return v.value
}

func (v *NullableExpandedTransaction) Set(val *ExpandedTransaction) {
	v.value = val
	v.isSet = true
}

func (v NullableExpandedTransaction) IsSet() bool {
	return v.isSet
}

func (v *NullableExpandedTransaction) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExpandedTransaction(val *ExpandedTransaction) *NullableExpandedTransaction {
	return &NullableExpandedTransaction{value: val, isSet: true}
}

func (v NullableExpandedTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExpandedTransaction) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


