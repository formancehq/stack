# WorkflowConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Stages** | **[]map[string]interface{}** |  | 

## Methods

### NewWorkflowConfig

`func NewWorkflowConfig(stages []map[string]interface{}, ) *WorkflowConfig`

NewWorkflowConfig instantiates a new WorkflowConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowConfigWithDefaults

`func NewWorkflowConfigWithDefaults() *WorkflowConfig`

NewWorkflowConfigWithDefaults instantiates a new WorkflowConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *WorkflowConfig) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowConfig) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowConfig) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *WorkflowConfig) HasName() bool`

HasName returns a boolean if a field has been set.

### GetStages

`func (o *WorkflowConfig) GetStages() []map[string]interface{}`

GetStages returns the Stages field if non-nil, zero value otherwise.

### GetStagesOk

`func (o *WorkflowConfig) GetStagesOk() (*[]map[string]interface{}, bool)`

GetStagesOk returns a tuple with the Stages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStages

`func (o *WorkflowConfig) SetStages(v []map[string]interface{})`

SetStages sets Stages field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


