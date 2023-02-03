# WorkflowInstanceHistory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Input** | [**Stage**](Stage.md) |  | 
**Error** | Pointer to **string** |  | [optional] 
**Terminated** | **bool** |  | 
**StartedAt** | **time.Time** |  | 
**TerminatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewWorkflowInstanceHistory

`func NewWorkflowInstanceHistory(name string, input Stage, terminated bool, startedAt time.Time, ) *WorkflowInstanceHistory`

NewWorkflowInstanceHistory instantiates a new WorkflowInstanceHistory object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowInstanceHistoryWithDefaults

`func NewWorkflowInstanceHistoryWithDefaults() *WorkflowInstanceHistory`

NewWorkflowInstanceHistoryWithDefaults instantiates a new WorkflowInstanceHistory object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *WorkflowInstanceHistory) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowInstanceHistory) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowInstanceHistory) SetName(v string)`

SetName sets Name field to given value.


### GetInput

`func (o *WorkflowInstanceHistory) GetInput() Stage`

GetInput returns the Input field if non-nil, zero value otherwise.

### GetInputOk

`func (o *WorkflowInstanceHistory) GetInputOk() (*Stage, bool)`

GetInputOk returns a tuple with the Input field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInput

`func (o *WorkflowInstanceHistory) SetInput(v Stage)`

SetInput sets Input field to given value.


### GetError

`func (o *WorkflowInstanceHistory) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *WorkflowInstanceHistory) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *WorkflowInstanceHistory) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *WorkflowInstanceHistory) HasError() bool`

HasError returns a boolean if a field has been set.

### GetTerminated

`func (o *WorkflowInstanceHistory) GetTerminated() bool`

GetTerminated returns the Terminated field if non-nil, zero value otherwise.

### GetTerminatedOk

`func (o *WorkflowInstanceHistory) GetTerminatedOk() (*bool, bool)`

GetTerminatedOk returns a tuple with the Terminated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminated

`func (o *WorkflowInstanceHistory) SetTerminated(v bool)`

SetTerminated sets Terminated field to given value.


### GetStartedAt

`func (o *WorkflowInstanceHistory) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *WorkflowInstanceHistory) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *WorkflowInstanceHistory) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.


### GetTerminatedAt

`func (o *WorkflowInstanceHistory) GetTerminatedAt() time.Time`

GetTerminatedAt returns the TerminatedAt field if non-nil, zero value otherwise.

### GetTerminatedAtOk

`func (o *WorkflowInstanceHistory) GetTerminatedAtOk() (*time.Time, bool)`

GetTerminatedAtOk returns a tuple with the TerminatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminatedAt

`func (o *WorkflowInstanceHistory) SetTerminatedAt(v time.Time)`

SetTerminatedAt sets TerminatedAt field to given value.

### HasTerminatedAt

`func (o *WorkflowInstanceHistory) HasTerminatedAt() bool`

HasTerminatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


