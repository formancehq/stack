# StageDelay

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Until** | Pointer to **time.Time** |  | [optional] 
**Duration** | Pointer to **string** |  | [optional] 

## Methods

### NewStageDelay

`func NewStageDelay() *StageDelay`

NewStageDelay instantiates a new StageDelay object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStageDelayWithDefaults

`func NewStageDelayWithDefaults() *StageDelay`

NewStageDelayWithDefaults instantiates a new StageDelay object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUntil

`func (o *StageDelay) GetUntil() time.Time`

GetUntil returns the Until field if non-nil, zero value otherwise.

### GetUntilOk

`func (o *StageDelay) GetUntilOk() (*time.Time, bool)`

GetUntilOk returns a tuple with the Until field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUntil

`func (o *StageDelay) SetUntil(v time.Time)`

SetUntil sets Until field to given value.

### HasUntil

`func (o *StageDelay) HasUntil() bool`

HasUntil returns a boolean if a field has been set.

### GetDuration

`func (o *StageDelay) GetDuration() string`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *StageDelay) GetDurationOk() (*string, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *StageDelay) SetDuration(v string)`

SetDuration sets Duration field to given value.

### HasDuration

`func (o *StageDelay) HasDuration() bool`

HasDuration returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


