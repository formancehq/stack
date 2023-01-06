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

// Payment struct for Payment
type Payment struct {
	Provider  string      `json:"provider"`
	Reference *string     `json:"reference,omitempty"`
	Scheme    string      `json:"scheme"`
	Status    string      `json:"status"`
	Type      string      `json:"type"`
	Id        string      `json:"id"`
	Amount    int32       `json:"amount"`
	Asset     string      `json:"asset"`
	Date      time.Time   `json:"date"`
	Raw       interface{} `json:"raw,omitempty"`
}

// NewPayment instantiates a new Payment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPayment(provider string, scheme string, status string, type_ string, id string, amount int32, asset string, date time.Time) *Payment {
	this := Payment{}
	this.Provider = provider
	this.Scheme = scheme
	this.Status = status
	this.Type = type_
	this.Id = id
	this.Amount = amount
	this.Asset = asset
	this.Date = date
	return &this
}

// NewPaymentWithDefaults instantiates a new Payment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaymentWithDefaults() *Payment {
	this := Payment{}
	return &this
}

// GetProvider returns the Provider field value
func (o *Payment) GetProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Provider
}

// GetProviderOk returns a tuple with the Provider field value
// and a boolean to check if the value has been set.
func (o *Payment) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Provider, true
}

// SetProvider sets field value
func (o *Payment) SetProvider(v string) {
	o.Provider = v
}

// GetReference returns the Reference field value if set, zero value otherwise.
func (o *Payment) GetReference() string {
	if o == nil || isNil(o.Reference) {
		var ret string
		return ret
	}
	return *o.Reference
}

// GetReferenceOk returns a tuple with the Reference field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Payment) GetReferenceOk() (*string, bool) {
	if o == nil || isNil(o.Reference) {
		return nil, false
	}
	return o.Reference, true
}

// HasReference returns a boolean if a field has been set.
func (o *Payment) HasReference() bool {
	if o != nil && !isNil(o.Reference) {
		return true
	}

	return false
}

// SetReference gets a reference to the given string and assigns it to the Reference field.
func (o *Payment) SetReference(v string) {
	o.Reference = &v
}

// GetScheme returns the Scheme field value
func (o *Payment) GetScheme() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Scheme
}

// GetSchemeOk returns a tuple with the Scheme field value
// and a boolean to check if the value has been set.
func (o *Payment) GetSchemeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Scheme, true
}

// SetScheme sets field value
func (o *Payment) SetScheme(v string) {
	o.Scheme = v
}

// GetStatus returns the Status field value
func (o *Payment) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *Payment) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *Payment) SetStatus(v string) {
	o.Status = v
}

// GetType returns the Type field value
func (o *Payment) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *Payment) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *Payment) SetType(v string) {
	o.Type = v
}

// GetId returns the Id field value
func (o *Payment) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Payment) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Payment) SetId(v string) {
	o.Id = v
}

// GetAmount returns the Amount field value
func (o *Payment) GetAmount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Amount
}

// GetAmountOk returns a tuple with the Amount field value
// and a boolean to check if the value has been set.
func (o *Payment) GetAmountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Amount, true
}

// SetAmount sets field value
func (o *Payment) SetAmount(v int32) {
	o.Amount = v
}

// GetAsset returns the Asset field value
func (o *Payment) GetAsset() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Asset
}

// GetAssetOk returns a tuple with the Asset field value
// and a boolean to check if the value has been set.
func (o *Payment) GetAssetOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Asset, true
}

// SetAsset sets field value
func (o *Payment) SetAsset(v string) {
	o.Asset = v
}

// GetDate returns the Date field value
func (o *Payment) GetDate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Date
}

// GetDateOk returns a tuple with the Date field value
// and a boolean to check if the value has been set.
func (o *Payment) GetDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Date, true
}

// SetDate sets field value
func (o *Payment) SetDate(v time.Time) {
	o.Date = v
}

// GetRaw returns the Raw field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Payment) GetRaw() interface{} {
	if o == nil {
		var ret interface{}
		return ret
	}
	return o.Raw
}

// GetRawOk returns a tuple with the Raw field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Payment) GetRawOk() (*interface{}, bool) {
	if o == nil || isNil(o.Raw) {
		return nil, false
	}
	return &o.Raw, true
}

// HasRaw returns a boolean if a field has been set.
func (o *Payment) HasRaw() bool {
	if o != nil && isNil(o.Raw) {
		return true
	}

	return false
}

// SetRaw gets a reference to the given interface{} and assigns it to the Raw field.
func (o *Payment) SetRaw(v interface{}) {
	o.Raw = v
}

func (o Payment) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["provider"] = o.Provider
	}
	if !isNil(o.Reference) {
		toSerialize["reference"] = o.Reference
	}
	if true {
		toSerialize["scheme"] = o.Scheme
	}
	if true {
		toSerialize["status"] = o.Status
	}
	if true {
		toSerialize["type"] = o.Type
	}
	if true {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["amount"] = o.Amount
	}
	if true {
		toSerialize["asset"] = o.Asset
	}
	if true {
		toSerialize["date"] = o.Date
	}
	if o.Raw != nil {
		toSerialize["raw"] = o.Raw
	}
	return json.Marshal(toSerialize)
}

type NullablePayment struct {
	value *Payment
	isSet bool
}

func (v NullablePayment) Get() *Payment {
	return v.value
}

func (v *NullablePayment) Set(val *Payment) {
	v.value = val
	v.isSet = true
}

func (v NullablePayment) IsSet() bool {
	return v.isSet
}

func (v *NullablePayment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePayment(val *Payment) *NullablePayment {
	return &NullablePayment{value: val, isSet: true}
}

func (v NullablePayment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePayment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
