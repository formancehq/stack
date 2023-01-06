# StripeTransferRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | Pointer to **interface{}** |  | [optional] 
**Asset** | Pointer to **interface{}** |  | [optional] 
**Destination** | Pointer to **interface{}** |  | [optional] 
**Metadata** | Pointer to **interface{}** | A set of key/value pairs that you can attach to a transfer object. It can be useful for storing additional information about the transfer in a structured format.  | [optional] 

## Methods

### NewStripeTransferRequest

`func NewStripeTransferRequest() *StripeTransferRequest`

NewStripeTransferRequest instantiates a new StripeTransferRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStripeTransferRequestWithDefaults

`func NewStripeTransferRequestWithDefaults() *StripeTransferRequest`

NewStripeTransferRequestWithDefaults instantiates a new StripeTransferRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *StripeTransferRequest) GetAmount() interface{}`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *StripeTransferRequest) GetAmountOk() (*interface{}, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *StripeTransferRequest) SetAmount(v interface{})`

SetAmount sets Amount field to given value.

### HasAmount

`func (o *StripeTransferRequest) HasAmount() bool`

HasAmount returns a boolean if a field has been set.

### SetAmountNil

`func (o *StripeTransferRequest) SetAmountNil(b bool)`

 SetAmountNil sets the value for Amount to be an explicit nil

### UnsetAmount
`func (o *StripeTransferRequest) UnsetAmount()`

UnsetAmount ensures that no value is present for Amount, not even an explicit nil
### GetAsset

`func (o *StripeTransferRequest) GetAsset() interface{}`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *StripeTransferRequest) GetAssetOk() (*interface{}, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *StripeTransferRequest) SetAsset(v interface{})`

SetAsset sets Asset field to given value.

### HasAsset

`func (o *StripeTransferRequest) HasAsset() bool`

HasAsset returns a boolean if a field has been set.

### SetAssetNil

`func (o *StripeTransferRequest) SetAssetNil(b bool)`

 SetAssetNil sets the value for Asset to be an explicit nil

### UnsetAsset
`func (o *StripeTransferRequest) UnsetAsset()`

UnsetAsset ensures that no value is present for Asset, not even an explicit nil
### GetDestination

`func (o *StripeTransferRequest) GetDestination() interface{}`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *StripeTransferRequest) GetDestinationOk() (*interface{}, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *StripeTransferRequest) SetDestination(v interface{})`

SetDestination sets Destination field to given value.

### HasDestination

`func (o *StripeTransferRequest) HasDestination() bool`

HasDestination returns a boolean if a field has been set.

### SetDestinationNil

`func (o *StripeTransferRequest) SetDestinationNil(b bool)`

 SetDestinationNil sets the value for Destination to be an explicit nil

### UnsetDestination
`func (o *StripeTransferRequest) UnsetDestination()`

UnsetDestination ensures that no value is present for Destination, not even an explicit nil
### GetMetadata

`func (o *StripeTransferRequest) GetMetadata() interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *StripeTransferRequest) GetMetadataOk() (*interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *StripeTransferRequest) SetMetadata(v interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *StripeTransferRequest) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### SetMetadataNil

`func (o *StripeTransferRequest) SetMetadataNil(b bool)`

 SetMetadataNil sets the value for Metadata to be an explicit nil

### UnsetMetadata
`func (o *StripeTransferRequest) UnsetMetadata()`

UnsetMetadata ensures that no value is present for Metadata, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


