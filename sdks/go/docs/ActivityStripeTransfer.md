# ActivityStripeTransfer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | Pointer to **int64** |  | [optional] 
**Asset** | Pointer to **string** |  | [optional] 
**Destination** | Pointer to **string** |  | [optional] 
**Metadata** | Pointer to **map[string]interface{}** | A set of key/value pairs that you can attach to a transfer object. It can be useful for storing additional information about the transfer in a structured format.  | [optional] 

## Methods

### NewActivityStripeTransfer

`func NewActivityStripeTransfer() *ActivityStripeTransfer`

NewActivityStripeTransfer instantiates a new ActivityStripeTransfer object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActivityStripeTransferWithDefaults

`func NewActivityStripeTransferWithDefaults() *ActivityStripeTransfer`

NewActivityStripeTransferWithDefaults instantiates a new ActivityStripeTransfer object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *ActivityStripeTransfer) GetAmount() int64`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *ActivityStripeTransfer) GetAmountOk() (*int64, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *ActivityStripeTransfer) SetAmount(v int64)`

SetAmount sets Amount field to given value.

### HasAmount

`func (o *ActivityStripeTransfer) HasAmount() bool`

HasAmount returns a boolean if a field has been set.

### GetAsset

`func (o *ActivityStripeTransfer) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *ActivityStripeTransfer) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *ActivityStripeTransfer) SetAsset(v string)`

SetAsset sets Asset field to given value.

### HasAsset

`func (o *ActivityStripeTransfer) HasAsset() bool`

HasAsset returns a boolean if a field has been set.

### GetDestination

`func (o *ActivityStripeTransfer) GetDestination() string`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *ActivityStripeTransfer) GetDestinationOk() (*string, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *ActivityStripeTransfer) SetDestination(v string)`

SetDestination sets Destination field to given value.

### HasDestination

`func (o *ActivityStripeTransfer) HasDestination() bool`

HasDestination returns a boolean if a field has been set.

### GetMetadata

`func (o *ActivityStripeTransfer) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ActivityStripeTransfer) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ActivityStripeTransfer) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ActivityStripeTransfer) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


