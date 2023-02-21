# WorkflowInstanceHistoryStage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Input** | [**WorkflowInstanceHistoryStageInput**](WorkflowInstanceHistoryStageInput.md) |  | 
**Output** | Pointer to [**WorkflowInstanceHistoryStageOutput**](WorkflowInstanceHistoryStageOutput.md) |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 
**Terminated** | **bool** |  | 
**StartedAt** | **time.Time** |  | 
**TerminatedAt** | Pointer to **time.Time** |  | [optional] 
**LastFailure** | Pointer to **string** |  | [optional] 
**Attempt** | **int32** |  | 
**NextExecution** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewWorkflowInstanceHistoryStage

`func NewWorkflowInstanceHistoryStage(name string, input WorkflowInstanceHistoryStageInput, terminated bool, startedAt time.Time, attempt int32, ) *WorkflowInstanceHistoryStage`

NewWorkflowInstanceHistoryStage instantiates a new WorkflowInstanceHistoryStage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowInstanceHistoryStageWithDefaults

`func NewWorkflowInstanceHistoryStageWithDefaults() *WorkflowInstanceHistoryStage`

NewWorkflowInstanceHistoryStageWithDefaults instantiates a new WorkflowInstanceHistoryStage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *WorkflowInstanceHistoryStage) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowInstanceHistoryStage) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowInstanceHistoryStage) SetName(v string)`

SetName sets Name field to given value.


### GetInput

`func (o *WorkflowInstanceHistoryStage) GetInput() WorkflowInstanceHistoryStageInput`

GetInput returns the Input field if non-nil, zero value otherwise.

### GetInputOk

`func (o *WorkflowInstanceHistoryStage) GetInputOk() (*WorkflowInstanceHistoryStageInput, bool)`

GetInputOk returns a tuple with the Input field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInput

`func (o *WorkflowInstanceHistoryStage) SetInput(v WorkflowInstanceHistoryStageInput)`

SetInput sets Input field to given value.


### GetOutput

`func (o *WorkflowInstanceHistoryStage) GetOutput() WorkflowInstanceHistoryStageOutput`

GetOutput returns the Output field if non-nil, zero value otherwise.

### GetOutputOk

`func (o *WorkflowInstanceHistoryStage) GetOutputOk() (*WorkflowInstanceHistoryStageOutput, bool)`

GetOutputOk returns a tuple with the Output field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutput

`func (o *WorkflowInstanceHistoryStage) SetOutput(v WorkflowInstanceHistoryStageOutput)`

SetOutput sets Output field to given value.

### HasOutput

`func (o *WorkflowInstanceHistoryStage) HasOutput() bool`

HasOutput returns a boolean if a field has been set.

### GetError

`func (o *WorkflowInstanceHistoryStage) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *WorkflowInstanceHistoryStage) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *WorkflowInstanceHistoryStage) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *WorkflowInstanceHistoryStage) HasError() bool`

HasError returns a boolean if a field has been set.

### GetTerminated

`func (o *WorkflowInstanceHistoryStage) GetTerminated() bool`

GetTerminated returns the Terminated field if non-nil, zero value otherwise.

### GetTerminatedOk

`func (o *WorkflowInstanceHistoryStage) GetTerminatedOk() (*bool, bool)`

GetTerminatedOk returns a tuple with the Terminated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminated

`func (o *WorkflowInstanceHistoryStage) SetTerminated(v bool)`

SetTerminated sets Terminated field to given value.


### GetStartedAt

`func (o *WorkflowInstanceHistoryStage) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *WorkflowInstanceHistoryStage) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *WorkflowInstanceHistoryStage) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.


### GetTerminatedAt

`func (o *WorkflowInstanceHistoryStage) GetTerminatedAt() time.Time`

GetTerminatedAt returns the TerminatedAt field if non-nil, zero value otherwise.

### GetTerminatedAtOk

`func (o *WorkflowInstanceHistoryStage) GetTerminatedAtOk() (*time.Time, bool)`

GetTerminatedAtOk returns a tuple with the TerminatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminatedAt

`func (o *WorkflowInstanceHistoryStage) SetTerminatedAt(v time.Time)`

SetTerminatedAt sets TerminatedAt field to given value.

### HasTerminatedAt

`func (o *WorkflowInstanceHistoryStage) HasTerminatedAt() bool`

HasTerminatedAt returns a boolean if a field has been set.

### GetLastFailure

`func (o *WorkflowInstanceHistoryStage) GetLastFailure() string`

GetLastFailure returns the LastFailure field if non-nil, zero value otherwise.

### GetLastFailureOk

`func (o *WorkflowInstanceHistoryStage) GetLastFailureOk() (*string, bool)`

GetLastFailureOk returns a tuple with the LastFailure field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastFailure

`func (o *WorkflowInstanceHistoryStage) SetLastFailure(v string)`

SetLastFailure sets LastFailure field to given value.

### HasLastFailure

`func (o *WorkflowInstanceHistoryStage) HasLastFailure() bool`

HasLastFailure returns a boolean if a field has been set.

### GetAttempt

`func (o *WorkflowInstanceHistoryStage) GetAttempt() int32`

GetAttempt returns the Attempt field if non-nil, zero value otherwise.

### GetAttemptOk

`func (o *WorkflowInstanceHistoryStage) GetAttemptOk() (*int32, bool)`

GetAttemptOk returns a tuple with the Attempt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttempt

`func (o *WorkflowInstanceHistoryStage) SetAttempt(v int32)`

SetAttempt sets Attempt field to given value.


### GetNextExecution

`func (o *WorkflowInstanceHistoryStage) GetNextExecution() time.Time`

GetNextExecution returns the NextExecution field if non-nil, zero value otherwise.

### GetNextExecutionOk

`func (o *WorkflowInstanceHistoryStage) GetNextExecutionOk() (*time.Time, bool)`

GetNextExecutionOk returns a tuple with the NextExecution field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextExecution

`func (o *WorkflowInstanceHistoryStage) SetNextExecution(v time.Time)`

SetNextExecution sets NextExecution field to given value.

### HasNextExecution

`func (o *WorkflowInstanceHistoryStage) HasNextExecution() bool`

HasNextExecution returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


