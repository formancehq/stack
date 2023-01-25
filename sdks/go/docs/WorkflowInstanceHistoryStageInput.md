# WorkflowInstanceHistoryStageInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Ledger** | **string** |  | 
**Data** | Pointer to [**DebitWalletRequest**](DebitWalletRequest.md) |  | [optional] 
**Amount** | Pointer to **int64** |  | [optional] 
**Asset** | Pointer to **string** |  | [optional] 
**Destination** | Pointer to **string** |  | [optional] 
**Metadata** | Pointer to **map[string]interface{}** | A set of key/value pairs that you can attach to a transfer object. It can be useful for storing additional information about the transfer in a structured format.  | [optional] 

## Methods

### NewWorkflowInstanceHistoryStageInput

`func NewWorkflowInstanceHistoryStageInput(id string, ledger string, ) *WorkflowInstanceHistoryStageInput`

NewWorkflowInstanceHistoryStageInput instantiates a new WorkflowInstanceHistoryStageInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowInstanceHistoryStageInputWithDefaults

`func NewWorkflowInstanceHistoryStageInputWithDefaults() *WorkflowInstanceHistoryStageInput`

NewWorkflowInstanceHistoryStageInputWithDefaults instantiates a new WorkflowInstanceHistoryStageInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *WorkflowInstanceHistoryStageInput) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WorkflowInstanceHistoryStageInput) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WorkflowInstanceHistoryStageInput) SetId(v string)`

SetId sets Id field to given value.


### GetLedger

`func (o *WorkflowInstanceHistoryStageInput) GetLedger() string`

GetLedger returns the Ledger field if non-nil, zero value otherwise.

### GetLedgerOk

`func (o *WorkflowInstanceHistoryStageInput) GetLedgerOk() (*string, bool)`

GetLedgerOk returns a tuple with the Ledger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedger

`func (o *WorkflowInstanceHistoryStageInput) SetLedger(v string)`

SetLedger sets Ledger field to given value.


### GetData

`func (o *WorkflowInstanceHistoryStageInput) GetData() DebitWalletRequest`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *WorkflowInstanceHistoryStageInput) GetDataOk() (*DebitWalletRequest, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *WorkflowInstanceHistoryStageInput) SetData(v DebitWalletRequest)`

SetData sets Data field to given value.

### HasData

`func (o *WorkflowInstanceHistoryStageInput) HasData() bool`

HasData returns a boolean if a field has been set.

### GetAmount

`func (o *WorkflowInstanceHistoryStageInput) GetAmount() int64`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *WorkflowInstanceHistoryStageInput) GetAmountOk() (*int64, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *WorkflowInstanceHistoryStageInput) SetAmount(v int64)`

SetAmount sets Amount field to given value.

### HasAmount

`func (o *WorkflowInstanceHistoryStageInput) HasAmount() bool`

HasAmount returns a boolean if a field has been set.

### GetAsset

`func (o *WorkflowInstanceHistoryStageInput) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *WorkflowInstanceHistoryStageInput) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *WorkflowInstanceHistoryStageInput) SetAsset(v string)`

SetAsset sets Asset field to given value.

### HasAsset

`func (o *WorkflowInstanceHistoryStageInput) HasAsset() bool`

HasAsset returns a boolean if a field has been set.

### GetDestination

`func (o *WorkflowInstanceHistoryStageInput) GetDestination() string`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *WorkflowInstanceHistoryStageInput) GetDestinationOk() (*string, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *WorkflowInstanceHistoryStageInput) SetDestination(v string)`

SetDestination sets Destination field to given value.

### HasDestination

`func (o *WorkflowInstanceHistoryStageInput) HasDestination() bool`

HasDestination returns a boolean if a field has been set.

### GetMetadata

`func (o *WorkflowInstanceHistoryStageInput) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *WorkflowInstanceHistoryStageInput) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *WorkflowInstanceHistoryStageInput) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *WorkflowInstanceHistoryStageInput) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


