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

// Cursor struct for Cursor
type Cursor struct {
	PageSize *int64 `json:"pageSize,omitempty"`
	HasMore *bool `json:"hasMore,omitempty"`
	Total *Total `json:"total,omitempty"`
	Next *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	Data []interface{} `json:"data,omitempty"`
}

// NewCursor instantiates a new Cursor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCursor() *Cursor {
	this := Cursor{}
	return &this
}

// NewCursorWithDefaults instantiates a new Cursor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCursorWithDefaults() *Cursor {
	this := Cursor{}
	return &this
}

// GetPageSize returns the PageSize field value if set, zero value otherwise.
func (o *Cursor) GetPageSize() int64 {
	if o == nil || isNil(o.PageSize) {
		var ret int64
		return ret
	}
	return *o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetPageSizeOk() (*int64, bool) {
	if o == nil || isNil(o.PageSize) {
    return nil, false
	}
	return o.PageSize, true
}

// HasPageSize returns a boolean if a field has been set.
func (o *Cursor) HasPageSize() bool {
	if o != nil && !isNil(o.PageSize) {
		return true
	}

	return false
}

// SetPageSize gets a reference to the given int64 and assigns it to the PageSize field.
func (o *Cursor) SetPageSize(v int64) {
	o.PageSize = &v
}

// GetHasMore returns the HasMore field value if set, zero value otherwise.
func (o *Cursor) GetHasMore() bool {
	if o == nil || isNil(o.HasMore) {
		var ret bool
		return ret
	}
	return *o.HasMore
}

// GetHasMoreOk returns a tuple with the HasMore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetHasMoreOk() (*bool, bool) {
	if o == nil || isNil(o.HasMore) {
    return nil, false
	}
	return o.HasMore, true
}

// HasHasMore returns a boolean if a field has been set.
func (o *Cursor) HasHasMore() bool {
	if o != nil && !isNil(o.HasMore) {
		return true
	}

	return false
}

// SetHasMore gets a reference to the given bool and assigns it to the HasMore field.
func (o *Cursor) SetHasMore(v bool) {
	o.HasMore = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *Cursor) GetTotal() Total {
	if o == nil || isNil(o.Total) {
		var ret Total
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetTotalOk() (*Total, bool) {
	if o == nil || isNil(o.Total) {
    return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *Cursor) HasTotal() bool {
	if o != nil && !isNil(o.Total) {
		return true
	}

	return false
}

// SetTotal gets a reference to the given Total and assigns it to the Total field.
func (o *Cursor) SetTotal(v Total) {
	o.Total = &v
}

// GetNext returns the Next field value if set, zero value otherwise.
func (o *Cursor) GetNext() string {
	if o == nil || isNil(o.Next) {
		var ret string
		return ret
	}
	return *o.Next
}

// GetNextOk returns a tuple with the Next field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetNextOk() (*string, bool) {
	if o == nil || isNil(o.Next) {
    return nil, false
	}
	return o.Next, true
}

// HasNext returns a boolean if a field has been set.
func (o *Cursor) HasNext() bool {
	if o != nil && !isNil(o.Next) {
		return true
	}

	return false
}

// SetNext gets a reference to the given string and assigns it to the Next field.
func (o *Cursor) SetNext(v string) {
	o.Next = &v
}

// GetPrevious returns the Previous field value if set, zero value otherwise.
func (o *Cursor) GetPrevious() string {
	if o == nil || isNil(o.Previous) {
		var ret string
		return ret
	}
	return *o.Previous
}

// GetPreviousOk returns a tuple with the Previous field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetPreviousOk() (*string, bool) {
	if o == nil || isNil(o.Previous) {
    return nil, false
	}
	return o.Previous, true
}

// HasPrevious returns a boolean if a field has been set.
func (o *Cursor) HasPrevious() bool {
	if o != nil && !isNil(o.Previous) {
		return true
	}

	return false
}

// SetPrevious gets a reference to the given string and assigns it to the Previous field.
func (o *Cursor) SetPrevious(v string) {
	o.Previous = &v
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *Cursor) GetData() []interface{} {
	if o == nil || isNil(o.Data) {
		var ret []interface{}
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Cursor) GetDataOk() ([]interface{}, bool) {
	if o == nil || isNil(o.Data) {
    return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *Cursor) HasData() bool {
	if o != nil && !isNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given []interface{} and assigns it to the Data field.
func (o *Cursor) SetData(v []interface{}) {
	o.Data = v
}

func (o Cursor) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.PageSize) {
		toSerialize["pageSize"] = o.PageSize
	}
	if !isNil(o.HasMore) {
		toSerialize["hasMore"] = o.HasMore
	}
	if !isNil(o.Total) {
		toSerialize["total"] = o.Total
	}
	if !isNil(o.Next) {
		toSerialize["next"] = o.Next
	}
	if !isNil(o.Previous) {
		toSerialize["previous"] = o.Previous
	}
	if !isNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableCursor struct {
	value *Cursor
	isSet bool
}

func (v NullableCursor) Get() *Cursor {
	return v.value
}

func (v *NullableCursor) Set(val *Cursor) {
	v.value = val
	v.isSet = true
}

func (v NullableCursor) IsSet() bool {
	return v.isSet
}

func (v *NullableCursor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCursor(val *Cursor) *NullableCursor {
	return &NullableCursor{value: val, isSet: true}
}

func (v NullableCursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCursor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


