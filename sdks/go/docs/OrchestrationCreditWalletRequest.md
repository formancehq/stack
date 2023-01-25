# OrchestrationCreditWalletRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | [**OrchestrationMonetary**](OrchestrationMonetary.md) |  | 
**Metadata** | Pointer to **map[string]interface{}** | Metadata associated with the wallet. | [optional] 
**Reference** | Pointer to **string** |  | [optional] 
**Sources** | [**[]Subject**](Subject.md) |  | 
**Balance** | Pointer to **string** | The balance to credit | [optional] 

## Methods

### NewOrchestrationCreditWalletRequest

`func NewOrchestrationCreditWalletRequest(amount OrchestrationMonetary, sources []Subject, ) *OrchestrationCreditWalletRequest`

NewOrchestrationCreditWalletRequest instantiates a new OrchestrationCreditWalletRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrchestrationCreditWalletRequestWithDefaults

`func NewOrchestrationCreditWalletRequestWithDefaults() *OrchestrationCreditWalletRequest`

NewOrchestrationCreditWalletRequestWithDefaults instantiates a new OrchestrationCreditWalletRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *OrchestrationCreditWalletRequest) GetAmount() OrchestrationMonetary`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *OrchestrationCreditWalletRequest) GetAmountOk() (*OrchestrationMonetary, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *OrchestrationCreditWalletRequest) SetAmount(v OrchestrationMonetary)`

SetAmount sets Amount field to given value.


### GetMetadata

`func (o *OrchestrationCreditWalletRequest) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *OrchestrationCreditWalletRequest) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *OrchestrationCreditWalletRequest) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *OrchestrationCreditWalletRequest) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetReference

`func (o *OrchestrationCreditWalletRequest) GetReference() string`

GetReference returns the Reference field if non-nil, zero value otherwise.

### GetReferenceOk

`func (o *OrchestrationCreditWalletRequest) GetReferenceOk() (*string, bool)`

GetReferenceOk returns a tuple with the Reference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReference

`func (o *OrchestrationCreditWalletRequest) SetReference(v string)`

SetReference sets Reference field to given value.

### HasReference

`func (o *OrchestrationCreditWalletRequest) HasReference() bool`

HasReference returns a boolean if a field has been set.

### GetSources

`func (o *OrchestrationCreditWalletRequest) GetSources() []Subject`

GetSources returns the Sources field if non-nil, zero value otherwise.

### GetSourcesOk

`func (o *OrchestrationCreditWalletRequest) GetSourcesOk() (*[]Subject, bool)`

GetSourcesOk returns a tuple with the Sources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSources

`func (o *OrchestrationCreditWalletRequest) SetSources(v []Subject)`

SetSources sets Sources field to given value.


### GetBalance

`func (o *OrchestrationCreditWalletRequest) GetBalance() string`

GetBalance returns the Balance field if non-nil, zero value otherwise.

### GetBalanceOk

`func (o *OrchestrationCreditWalletRequest) GetBalanceOk() (*string, bool)`

GetBalanceOk returns a tuple with the Balance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalance

`func (o *OrchestrationCreditWalletRequest) SetBalance(v string)`

SetBalance sets Balance field to given value.

### HasBalance

`func (o *OrchestrationCreditWalletRequest) HasBalance() bool`

HasBalance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


