# Payment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Provider** | **interface{}** |  |
**Reference** | Pointer to **interface{}** |  | [optional]
**Scheme** | **interface{}** |  |
**Status** | **interface{}** |  |
**Type** | **interface{}** |  |
**Id** | **interface{}** |  |
**Amount** | **interface{}** |  |
**Asset** | **interface{}** |  |
**Date** | **interface{}** |  |
**Raw** | Pointer to **interface{}** |  | [optional]

## Methods

### NewPayment

`func NewPayment(provider interface{}, scheme interface{}, status interface{}, type_ interface{}, id interface{}, amount interface{}, asset interface{}, date interface{}, ) *Payment`

NewPayment instantiates a new Payment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaymentWithDefaults

`func NewPaymentWithDefaults() *Payment`

NewPaymentWithDefaults instantiates a new Payment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProvider

`func (o *Payment) GetProvider() interface{}`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *Payment) GetProviderOk() (*interface{}, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *Payment) SetProvider(v interface{})`

SetProvider sets Provider field to given value.


### SetProviderNil

`func (o *Payment) SetProviderNil(b bool)`

 SetProviderNil sets the value for Provider to be an explicit nil

### UnsetProvider
`func (o *Payment) UnsetProvider()`

UnsetProvider ensures that no value is present for Provider, not even an explicit nil
### GetReference

`func (o *Payment) GetReference() interface{}`

GetReference returns the Reference field if non-nil, zero value otherwise.

### GetReferenceOk

`func (o *Payment) GetReferenceOk() (*interface{}, bool)`

GetReferenceOk returns a tuple with the Reference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReference

`func (o *Payment) SetReference(v interface{})`

SetReference sets Reference field to given value.

### HasReference

`func (o *Payment) HasReference() bool`

HasReference returns a boolean if a field has been set.

### SetReferenceNil

`func (o *Payment) SetReferenceNil(b bool)`

 SetReferenceNil sets the value for Reference to be an explicit nil

### UnsetReference
`func (o *Payment) UnsetReference()`

UnsetReference ensures that no value is present for Reference, not even an explicit nil
### GetScheme

`func (o *Payment) GetScheme() interface{}`

GetScheme returns the Scheme field if non-nil, zero value otherwise.

### GetSchemeOk

`func (o *Payment) GetSchemeOk() (*interface{}, bool)`

GetSchemeOk returns a tuple with the Scheme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheme

`func (o *Payment) SetScheme(v interface{})`

SetScheme sets Scheme field to given value.


### SetSchemeNil

`func (o *Payment) SetSchemeNil(b bool)`

 SetSchemeNil sets the value for Scheme to be an explicit nil

### UnsetScheme
`func (o *Payment) UnsetScheme()`

UnsetScheme ensures that no value is present for Scheme, not even an explicit nil
### GetStatus

`func (o *Payment) GetStatus() interface{}`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Payment) GetStatusOk() (*interface{}, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Payment) SetStatus(v interface{})`

SetStatus sets Status field to given value.


### SetStatusNil

`func (o *Payment) SetStatusNil(b bool)`

 SetStatusNil sets the value for Status to be an explicit nil

### UnsetStatus
`func (o *Payment) UnsetStatus()`

UnsetStatus ensures that no value is present for Status, not even an explicit nil
### GetType

`func (o *Payment) GetType() interface{}`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Payment) GetTypeOk() (*interface{}, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Payment) SetType(v interface{})`

SetType sets Type field to given value.


### SetTypeNil

`func (o *Payment) SetTypeNil(b bool)`

 SetTypeNil sets the value for Type to be an explicit nil

### UnsetType
`func (o *Payment) UnsetType()`

UnsetType ensures that no value is present for Type, not even an explicit nil
### GetId

`func (o *Payment) GetId() interface{}`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Payment) GetIdOk() (*interface{}, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Payment) SetId(v interface{})`

SetId sets Id field to given value.


### SetIdNil

`func (o *Payment) SetIdNil(b bool)`

 SetIdNil sets the value for Id to be an explicit nil

### UnsetId
`func (o *Payment) UnsetId()`

UnsetId ensures that no value is present for Id, not even an explicit nil
### GetAmount

`func (o *Payment) GetAmount() interface{}`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *Payment) GetAmountOk() (*interface{}, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *Payment) SetAmount(v interface{})`

SetAmount sets Amount field to given value.


### SetAmountNil

`func (o *Payment) SetAmountNil(b bool)`

 SetAmountNil sets the value for Amount to be an explicit nil

### UnsetAmount
`func (o *Payment) UnsetAmount()`

UnsetAmount ensures that no value is present for Amount, not even an explicit nil
### GetAsset

`func (o *Payment) GetAsset() interface{}`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *Payment) GetAssetOk() (*interface{}, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *Payment) SetAsset(v interface{})`

SetAsset sets Asset field to given value.


### SetAssetNil

`func (o *Payment) SetAssetNil(b bool)`

 SetAssetNil sets the value for Asset to be an explicit nil

### UnsetAsset
`func (o *Payment) UnsetAsset()`

UnsetAsset ensures that no value is present for Asset, not even an explicit nil
### GetDate

`func (o *Payment) GetDate() interface{}`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *Payment) GetDateOk() (*interface{}, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *Payment) SetDate(v interface{})`

SetDate sets Date field to given value.


### SetDateNil

`func (o *Payment) SetDateNil(b bool)`

 SetDateNil sets the value for Date to be an explicit nil

### UnsetDate
`func (o *Payment) UnsetDate()`

UnsetDate ensures that no value is present for Date, not even an explicit nil
### GetRaw

`func (o *Payment) GetRaw() interface{}`

GetRaw returns the Raw field if non-nil, zero value otherwise.

### GetRawOk

`func (o *Payment) GetRawOk() (*interface{}, bool)`

GetRawOk returns a tuple with the Raw field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRaw

`func (o *Payment) SetRaw(v interface{})`

SetRaw sets Raw field to given value.

### HasRaw

`func (o *Payment) HasRaw() bool`

HasRaw returns a boolean if a field has been set.

### SetRawNil

`func (o *Payment) SetRawNil(b bool)`

 SetRawNil sets the value for Raw to be an explicit nil

### UnsetRaw
`func (o *Payment) UnsetRaw()`

UnsetRaw ensures that no value is present for Raw, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
