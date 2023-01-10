# TaskWise

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ConnectorID** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**Descriptor** | Pointer to [**TaskWiseDescriptor**](TaskWiseDescriptor.md) |  | [optional] 
**Status** | Pointer to [**PaymentStatus**](PaymentStatus.md) |  | [optional] 
**State** | Pointer to **map[string]interface{}** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewTaskWise

`func NewTaskWise() *TaskWise`

NewTaskWise instantiates a new TaskWise object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaskWiseWithDefaults

`func NewTaskWiseWithDefaults() *TaskWise`

NewTaskWiseWithDefaults instantiates a new TaskWise object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TaskWise) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaskWise) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaskWise) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TaskWise) HasId() bool`

HasId returns a boolean if a field has been set.

### GetConnectorID

`func (o *TaskWise) GetConnectorID() string`

GetConnectorID returns the ConnectorID field if non-nil, zero value otherwise.

### GetConnectorIDOk

`func (o *TaskWise) GetConnectorIDOk() (*string, bool)`

GetConnectorIDOk returns a tuple with the ConnectorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorID

`func (o *TaskWise) SetConnectorID(v string)`

SetConnectorID sets ConnectorID field to given value.

### HasConnectorID

`func (o *TaskWise) HasConnectorID() bool`

HasConnectorID returns a boolean if a field has been set.

### GetCreatedAt

`func (o *TaskWise) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *TaskWise) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *TaskWise) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *TaskWise) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *TaskWise) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *TaskWise) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *TaskWise) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *TaskWise) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetDescriptor

`func (o *TaskWise) GetDescriptor() TaskWiseDescriptor`

GetDescriptor returns the Descriptor field if non-nil, zero value otherwise.

### GetDescriptorOk

`func (o *TaskWise) GetDescriptorOk() (*TaskWiseDescriptor, bool)`

GetDescriptorOk returns a tuple with the Descriptor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptor

`func (o *TaskWise) SetDescriptor(v TaskWiseDescriptor)`

SetDescriptor sets Descriptor field to given value.

### HasDescriptor

`func (o *TaskWise) HasDescriptor() bool`

HasDescriptor returns a boolean if a field has been set.

### GetStatus

`func (o *TaskWise) GetStatus() PaymentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *TaskWise) GetStatusOk() (*PaymentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *TaskWise) SetStatus(v PaymentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *TaskWise) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetState

`func (o *TaskWise) GetState() map[string]interface{}`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *TaskWise) GetStateOk() (*map[string]interface{}, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *TaskWise) SetState(v map[string]interface{})`

SetState sets State field to given value.

### HasState

`func (o *TaskWise) HasState() bool`

HasState returns a boolean if a field has been set.

### GetError

`func (o *TaskWise) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *TaskWise) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *TaskWise) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *TaskWise) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


