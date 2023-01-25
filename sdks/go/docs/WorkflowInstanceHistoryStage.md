# WorkflowInstanceHistoryStage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Input** | [**WorkflowInstanceHistoryStageInput**](WorkflowInstanceHistoryStageInput.md) |  | 
**Output** | Pointer to [**WorkflowInstanceHistoryStageOutput**](WorkflowInstanceHistoryStageOutput.md) |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkflowInstanceHistoryStage

`func NewWorkflowInstanceHistoryStage(name string, input WorkflowInstanceHistoryStageInput, ) *WorkflowInstanceHistoryStage`

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


