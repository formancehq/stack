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

// LogsCursorResponseCursor struct for LogsCursorResponseCursor
type LogsCursorResponseCursor struct {
	PageSize int64 `json:"pageSize"`
	HasMore bool `json:"hasMore"`
	Previous *string `json:"previous,omitempty"`
	Next *string `json:"next,omitempty"`
	Data []Log `json:"data"`
}

// NewLogsCursorResponseCursor instantiates a new LogsCursorResponseCursor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLogsCursorResponseCursor(pageSize int64, hasMore bool, data []Log) *LogsCursorResponseCursor {
	this := LogsCursorResponseCursor{}
	this.PageSize = pageSize
	this.HasMore = hasMore
	this.Data = data
	return &this
}

// NewLogsCursorResponseCursorWithDefaults instantiates a new LogsCursorResponseCursor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLogsCursorResponseCursorWithDefaults() *LogsCursorResponseCursor {
	this := LogsCursorResponseCursor{}
	return &this
}

// GetPageSize returns the PageSize field value
func (o *LogsCursorResponseCursor) GetPageSize() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value
// and a boolean to check if the value has been set.
func (o *LogsCursorResponseCursor) GetPageSizeOk() (*int64, bool) {
	if o == nil {
    return nil, false
	}
	return &o.PageSize, true
}

// SetPageSize sets field value
func (o *LogsCursorResponseCursor) SetPageSize(v int64) {
	o.PageSize = v
}

// GetHasMore returns the HasMore field value
func (o *LogsCursorResponseCursor) GetHasMore() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.HasMore
}

// GetHasMoreOk returns a tuple with the HasMore field value
// and a boolean to check if the value has been set.
func (o *LogsCursorResponseCursor) GetHasMoreOk() (*bool, bool) {
	if o == nil {
    return nil, false
	}
	return &o.HasMore, true
}

// SetHasMore sets field value
func (o *LogsCursorResponseCursor) SetHasMore(v bool) {
	o.HasMore = v
}

// GetPrevious returns the Previous field value if set, zero value otherwise.
func (o *LogsCursorResponseCursor) GetPrevious() string {
	if o == nil || isNil(o.Previous) {
		var ret string
		return ret
	}
	return *o.Previous
}

// GetPreviousOk returns a tuple with the Previous field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogsCursorResponseCursor) GetPreviousOk() (*string, bool) {
	if o == nil || isNil(o.Previous) {
    return nil, false
	}
	return o.Previous, true
}

// HasPrevious returns a boolean if a field has been set.
func (o *LogsCursorResponseCursor) HasPrevious() bool {
	if o != nil && !isNil(o.Previous) {
		return true
	}

	return false
}

// SetPrevious gets a reference to the given string and assigns it to the Previous field.
func (o *LogsCursorResponseCursor) SetPrevious(v string) {
	o.Previous = &v
}

// GetNext returns the Next field value if set, zero value otherwise.
func (o *LogsCursorResponseCursor) GetNext() string {
	if o == nil || isNil(o.Next) {
		var ret string
		return ret
	}
	return *o.Next
}

// GetNextOk returns a tuple with the Next field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogsCursorResponseCursor) GetNextOk() (*string, bool) {
	if o == nil || isNil(o.Next) {
    return nil, false
	}
	return o.Next, true
}

// HasNext returns a boolean if a field has been set.
func (o *LogsCursorResponseCursor) HasNext() bool {
	if o != nil && !isNil(o.Next) {
		return true
	}

	return false
}

// SetNext gets a reference to the given string and assigns it to the Next field.
func (o *LogsCursorResponseCursor) SetNext(v string) {
	o.Next = &v
}

// GetData returns the Data field value
func (o *LogsCursorResponseCursor) GetData() []Log {
	if o == nil {
		var ret []Log
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *LogsCursorResponseCursor) GetDataOk() ([]Log, bool) {
	if o == nil {
    return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *LogsCursorResponseCursor) SetData(v []Log) {
	o.Data = v
}

func (o LogsCursorResponseCursor) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["pageSize"] = o.PageSize
	}
	if true {
		toSerialize["hasMore"] = o.HasMore
	}
	if !isNil(o.Previous) {
		toSerialize["previous"] = o.Previous
	}
	if !isNil(o.Next) {
		toSerialize["next"] = o.Next
	}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableLogsCursorResponseCursor struct {
	value *LogsCursorResponseCursor
	isSet bool
}

func (v NullableLogsCursorResponseCursor) Get() *LogsCursorResponseCursor {
	return v.value
}

func (v *NullableLogsCursorResponseCursor) Set(val *LogsCursorResponseCursor) {
	v.value = val
	v.isSet = true
}

func (v NullableLogsCursorResponseCursor) IsSet() bool {
	return v.isSet
}

func (v *NullableLogsCursorResponseCursor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLogsCursorResponseCursor(val *LogsCursorResponseCursor) *NullableLogsCursorResponseCursor {
	return &NullableLogsCursorResponseCursor{value: val, isSet: true}
}

func (v NullableLogsCursorResponseCursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLogsCursorResponseCursor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


