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

// Query struct for Query
type Query struct {
	Ledgers []string `json:"ledgers,omitempty"`
	After []string `json:"after,omitempty"`
	PageSize *int64 `json:"pageSize,omitempty"`
	Terms []string `json:"terms,omitempty"`
	Sort *string `json:"sort,omitempty"`
	Policy *string `json:"policy,omitempty"`
	Target *string `json:"target,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}

// NewQuery instantiates a new Query object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQuery() *Query {
	this := Query{}
	return &this
}

// NewQueryWithDefaults instantiates a new Query object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQueryWithDefaults() *Query {
	this := Query{}
	return &this
}

// GetLedgers returns the Ledgers field value if set, zero value otherwise.
func (o *Query) GetLedgers() []string {
	if o == nil || isNil(o.Ledgers) {
		var ret []string
		return ret
	}
	return o.Ledgers
}

// GetLedgersOk returns a tuple with the Ledgers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetLedgersOk() ([]string, bool) {
	if o == nil || isNil(o.Ledgers) {
    return nil, false
	}
	return o.Ledgers, true
}

// HasLedgers returns a boolean if a field has been set.
func (o *Query) HasLedgers() bool {
	if o != nil && !isNil(o.Ledgers) {
		return true
	}

	return false
}

// SetLedgers gets a reference to the given []string and assigns it to the Ledgers field.
func (o *Query) SetLedgers(v []string) {
	o.Ledgers = v
}

// GetAfter returns the After field value if set, zero value otherwise.
func (o *Query) GetAfter() []string {
	if o == nil || isNil(o.After) {
		var ret []string
		return ret
	}
	return o.After
}

// GetAfterOk returns a tuple with the After field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetAfterOk() ([]string, bool) {
	if o == nil || isNil(o.After) {
    return nil, false
	}
	return o.After, true
}

// HasAfter returns a boolean if a field has been set.
func (o *Query) HasAfter() bool {
	if o != nil && !isNil(o.After) {
		return true
	}

	return false
}

// SetAfter gets a reference to the given []string and assigns it to the After field.
func (o *Query) SetAfter(v []string) {
	o.After = v
}

// GetPageSize returns the PageSize field value if set, zero value otherwise.
func (o *Query) GetPageSize() int64 {
	if o == nil || isNil(o.PageSize) {
		var ret int64
		return ret
	}
	return *o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetPageSizeOk() (*int64, bool) {
	if o == nil || isNil(o.PageSize) {
    return nil, false
	}
	return o.PageSize, true
}

// HasPageSize returns a boolean if a field has been set.
func (o *Query) HasPageSize() bool {
	if o != nil && !isNil(o.PageSize) {
		return true
	}

	return false
}

// SetPageSize gets a reference to the given int64 and assigns it to the PageSize field.
func (o *Query) SetPageSize(v int64) {
	o.PageSize = &v
}

// GetTerms returns the Terms field value if set, zero value otherwise.
func (o *Query) GetTerms() []string {
	if o == nil || isNil(o.Terms) {
		var ret []string
		return ret
	}
	return o.Terms
}

// GetTermsOk returns a tuple with the Terms field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetTermsOk() ([]string, bool) {
	if o == nil || isNil(o.Terms) {
    return nil, false
	}
	return o.Terms, true
}

// HasTerms returns a boolean if a field has been set.
func (o *Query) HasTerms() bool {
	if o != nil && !isNil(o.Terms) {
		return true
	}

	return false
}

// SetTerms gets a reference to the given []string and assigns it to the Terms field.
func (o *Query) SetTerms(v []string) {
	o.Terms = v
}

// GetSort returns the Sort field value if set, zero value otherwise.
func (o *Query) GetSort() string {
	if o == nil || isNil(o.Sort) {
		var ret string
		return ret
	}
	return *o.Sort
}

// GetSortOk returns a tuple with the Sort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetSortOk() (*string, bool) {
	if o == nil || isNil(o.Sort) {
    return nil, false
	}
	return o.Sort, true
}

// HasSort returns a boolean if a field has been set.
func (o *Query) HasSort() bool {
	if o != nil && !isNil(o.Sort) {
		return true
	}

	return false
}

// SetSort gets a reference to the given string and assigns it to the Sort field.
func (o *Query) SetSort(v string) {
	o.Sort = &v
}

// GetPolicy returns the Policy field value if set, zero value otherwise.
func (o *Query) GetPolicy() string {
	if o == nil || isNil(o.Policy) {
		var ret string
		return ret
	}
	return *o.Policy
}

// GetPolicyOk returns a tuple with the Policy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetPolicyOk() (*string, bool) {
	if o == nil || isNil(o.Policy) {
    return nil, false
	}
	return o.Policy, true
}

// HasPolicy returns a boolean if a field has been set.
func (o *Query) HasPolicy() bool {
	if o != nil && !isNil(o.Policy) {
		return true
	}

	return false
}

// SetPolicy gets a reference to the given string and assigns it to the Policy field.
func (o *Query) SetPolicy(v string) {
	o.Policy = &v
}

// GetTarget returns the Target field value if set, zero value otherwise.
func (o *Query) GetTarget() string {
	if o == nil || isNil(o.Target) {
		var ret string
		return ret
	}
	return *o.Target
}

// GetTargetOk returns a tuple with the Target field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetTargetOk() (*string, bool) {
	if o == nil || isNil(o.Target) {
    return nil, false
	}
	return o.Target, true
}

// HasTarget returns a boolean if a field has been set.
func (o *Query) HasTarget() bool {
	if o != nil && !isNil(o.Target) {
		return true
	}

	return false
}

// SetTarget gets a reference to the given string and assigns it to the Target field.
func (o *Query) SetTarget(v string) {
	o.Target = &v
}

// GetCursor returns the Cursor field value if set, zero value otherwise.
func (o *Query) GetCursor() string {
	if o == nil || isNil(o.Cursor) {
		var ret string
		return ret
	}
	return *o.Cursor
}

// GetCursorOk returns a tuple with the Cursor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Query) GetCursorOk() (*string, bool) {
	if o == nil || isNil(o.Cursor) {
    return nil, false
	}
	return o.Cursor, true
}

// HasCursor returns a boolean if a field has been set.
func (o *Query) HasCursor() bool {
	if o != nil && !isNil(o.Cursor) {
		return true
	}

	return false
}

// SetCursor gets a reference to the given string and assigns it to the Cursor field.
func (o *Query) SetCursor(v string) {
	o.Cursor = &v
}

func (o Query) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Ledgers) {
		toSerialize["ledgers"] = o.Ledgers
	}
	if !isNil(o.After) {
		toSerialize["after"] = o.After
	}
	if !isNil(o.PageSize) {
		toSerialize["pageSize"] = o.PageSize
	}
	if !isNil(o.Terms) {
		toSerialize["terms"] = o.Terms
	}
	if !isNil(o.Sort) {
		toSerialize["sort"] = o.Sort
	}
	if !isNil(o.Policy) {
		toSerialize["policy"] = o.Policy
	}
	if !isNil(o.Target) {
		toSerialize["target"] = o.Target
	}
	if !isNil(o.Cursor) {
		toSerialize["cursor"] = o.Cursor
	}
	return json.Marshal(toSerialize)
}

type NullableQuery struct {
	value *Query
	isSet bool
}

func (v NullableQuery) Get() *Query {
	return v.value
}

func (v *NullableQuery) Set(val *Query) {
	v.value = val
	v.isSet = true
}

func (v NullableQuery) IsSet() bool {
	return v.isSet
}

func (v *NullableQuery) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQuery(val *Query) *NullableQuery {
	return &NullableQuery{value: val, isSet: true}
}

func (v NullableQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQuery) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
