# Stage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | Pointer to [**Monetary**](Monetary.md) |  | [optional] 
**Destination** | Pointer to [**StageSendDestination**](StageSendDestination.md) |  | [optional] 
**Source** | Pointer to [**StageSendSource**](StageSendSource.md) |  | [optional] 

## Methods

### NewStage

`func NewStage() *Stage`

NewStage instantiates a new Stage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStageWithDefaults

`func NewStageWithDefaults() *Stage`

NewStageWithDefaults instantiates a new Stage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *Stage) GetAmount() Monetary`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *Stage) GetAmountOk() (*Monetary, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *Stage) SetAmount(v Monetary)`

SetAmount sets Amount field to given value.

### HasAmount

`func (o *Stage) HasAmount() bool`

HasAmount returns a boolean if a field has been set.

### GetDestination

`func (o *Stage) GetDestination() StageSendDestination`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *Stage) GetDestinationOk() (*StageSendDestination, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *Stage) SetDestination(v StageSendDestination)`

SetDestination sets Destination field to given value.

### HasDestination

`func (o *Stage) HasDestination() bool`

HasDestination returns a boolean if a field has been set.

### GetSource

`func (o *Stage) GetSource() StageSendSource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *Stage) GetSourceOk() (*StageSendSource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *Stage) SetSource(v StageSendSource)`

SetSource sets Source field to given value.

### HasSource

`func (o *Stage) HasSource() bool`

HasSource returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


