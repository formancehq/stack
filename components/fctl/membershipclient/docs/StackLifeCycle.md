# StackLifeCycle

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** |  | 
**State** | **string** |  | 
**ExpectedStatus** | **string** |  | 
**LastStateUpdate** | **time.Time** |  | 
**LastExpectedStatusUpdate** | **time.Time** |  | 
**LastStatusUpdate** | **time.Time** |  | 

## Methods

### NewStackLifeCycle

`func NewStackLifeCycle(status string, state string, expectedStatus string, lastStateUpdate time.Time, lastExpectedStatusUpdate time.Time, lastStatusUpdate time.Time, ) *StackLifeCycle`

NewStackLifeCycle instantiates a new StackLifeCycle object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStackLifeCycleWithDefaults

`func NewStackLifeCycleWithDefaults() *StackLifeCycle`

NewStackLifeCycleWithDefaults instantiates a new StackLifeCycle object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *StackLifeCycle) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *StackLifeCycle) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *StackLifeCycle) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetState

`func (o *StackLifeCycle) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *StackLifeCycle) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *StackLifeCycle) SetState(v string)`

SetState sets State field to given value.


### GetExpectedStatus

`func (o *StackLifeCycle) GetExpectedStatus() string`

GetExpectedStatus returns the ExpectedStatus field if non-nil, zero value otherwise.

### GetExpectedStatusOk

`func (o *StackLifeCycle) GetExpectedStatusOk() (*string, bool)`

GetExpectedStatusOk returns a tuple with the ExpectedStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedStatus

`func (o *StackLifeCycle) SetExpectedStatus(v string)`

SetExpectedStatus sets ExpectedStatus field to given value.


### GetLastStateUpdate

`func (o *StackLifeCycle) GetLastStateUpdate() time.Time`

GetLastStateUpdate returns the LastStateUpdate field if non-nil, zero value otherwise.

### GetLastStateUpdateOk

`func (o *StackLifeCycle) GetLastStateUpdateOk() (*time.Time, bool)`

GetLastStateUpdateOk returns a tuple with the LastStateUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStateUpdate

`func (o *StackLifeCycle) SetLastStateUpdate(v time.Time)`

SetLastStateUpdate sets LastStateUpdate field to given value.


### GetLastExpectedStatusUpdate

`func (o *StackLifeCycle) GetLastExpectedStatusUpdate() time.Time`

GetLastExpectedStatusUpdate returns the LastExpectedStatusUpdate field if non-nil, zero value otherwise.

### GetLastExpectedStatusUpdateOk

`func (o *StackLifeCycle) GetLastExpectedStatusUpdateOk() (*time.Time, bool)`

GetLastExpectedStatusUpdateOk returns a tuple with the LastExpectedStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastExpectedStatusUpdate

`func (o *StackLifeCycle) SetLastExpectedStatusUpdate(v time.Time)`

SetLastExpectedStatusUpdate sets LastExpectedStatusUpdate field to given value.


### GetLastStatusUpdate

`func (o *StackLifeCycle) GetLastStatusUpdate() time.Time`

GetLastStatusUpdate returns the LastStatusUpdate field if non-nil, zero value otherwise.

### GetLastStatusUpdateOk

`func (o *StackLifeCycle) GetLastStatusUpdateOk() (*time.Time, bool)`

GetLastStatusUpdateOk returns a tuple with the LastStatusUpdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStatusUpdate

`func (o *StackLifeCycle) SetLastStatusUpdate(v time.Time)`

SetLastStatusUpdate sets LastStatusUpdate field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


