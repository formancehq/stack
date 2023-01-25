# WorkflowInstanceHistory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Input** | [**Stage**](Stage.md) |  | 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkflowInstanceHistory

`func NewWorkflowInstanceHistory(name string, input Stage, ) *WorkflowInstanceHistory`

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


