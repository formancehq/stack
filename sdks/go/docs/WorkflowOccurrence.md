# WorkflowOccurrence

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowID** | **string** |  | 
**Id** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**UpdatedAt** | **time.Time** |  | 
**Status** | Pointer to [**[]StageStatus**](StageStatus.md) |  | [optional] 

## Methods

### NewWorkflowOccurrence

`func NewWorkflowOccurrence(workflowID string, id string, createdAt time.Time, updatedAt time.Time, ) *WorkflowOccurrence`

NewWorkflowOccurrence instantiates a new WorkflowOccurrence object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowOccurrenceWithDefaults

`func NewWorkflowOccurrenceWithDefaults() *WorkflowOccurrence`

NewWorkflowOccurrenceWithDefaults instantiates a new WorkflowOccurrence object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowID

`func (o *WorkflowOccurrence) GetWorkflowID() string`

GetWorkflowID returns the WorkflowID field if non-nil, zero value otherwise.

### GetWorkflowIDOk

`func (o *WorkflowOccurrence) GetWorkflowIDOk() (*string, bool)`

GetWorkflowIDOk returns a tuple with the WorkflowID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowID

`func (o *WorkflowOccurrence) SetWorkflowID(v string)`

SetWorkflowID sets WorkflowID field to given value.


### GetId

`func (o *WorkflowOccurrence) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WorkflowOccurrence) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WorkflowOccurrence) SetId(v string)`

SetId sets Id field to given value.


### GetCreatedAt

`func (o *WorkflowOccurrence) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *WorkflowOccurrence) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *WorkflowOccurrence) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *WorkflowOccurrence) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *WorkflowOccurrence) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *WorkflowOccurrence) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetStatus

`func (o *WorkflowOccurrence) GetStatus() []StageStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WorkflowOccurrence) GetStatusOk() (*[]StageStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WorkflowOccurrence) SetStatus(v []StageStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *WorkflowOccurrence) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


