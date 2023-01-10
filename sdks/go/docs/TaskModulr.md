# TaskModulr

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ConnectorID** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**Descriptor** | Pointer to [**TaskModulrDescriptor**](TaskModulrDescriptor.md) |  | [optional] 
**Status** | Pointer to [**PaymentStatus**](PaymentStatus.md) |  | [optional] 
**State** | Pointer to **map[string]interface{}** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewTaskModulr

`func NewTaskModulr() *TaskModulr`

NewTaskModulr instantiates a new TaskModulr object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaskModulrWithDefaults

`func NewTaskModulrWithDefaults() *TaskModulr`

NewTaskModulrWithDefaults instantiates a new TaskModulr object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TaskModulr) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaskModulr) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaskModulr) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TaskModulr) HasId() bool`

HasId returns a boolean if a field has been set.

### GetConnectorID

`func (o *TaskModulr) GetConnectorID() string`

GetConnectorID returns the ConnectorID field if non-nil, zero value otherwise.

### GetConnectorIDOk

`func (o *TaskModulr) GetConnectorIDOk() (*string, bool)`

GetConnectorIDOk returns a tuple with the ConnectorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorID

`func (o *TaskModulr) SetConnectorID(v string)`

SetConnectorID sets ConnectorID field to given value.

### HasConnectorID

`func (o *TaskModulr) HasConnectorID() bool`

HasConnectorID returns a boolean if a field has been set.

### GetCreatedAt

`func (o *TaskModulr) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *TaskModulr) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *TaskModulr) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *TaskModulr) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *TaskModulr) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *TaskModulr) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *TaskModulr) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *TaskModulr) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetDescriptor

`func (o *TaskModulr) GetDescriptor() TaskModulrDescriptor`

GetDescriptor returns the Descriptor field if non-nil, zero value otherwise.

### GetDescriptorOk

`func (o *TaskModulr) GetDescriptorOk() (*TaskModulrDescriptor, bool)`

GetDescriptorOk returns a tuple with the Descriptor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptor

`func (o *TaskModulr) SetDescriptor(v TaskModulrDescriptor)`

SetDescriptor sets Descriptor field to given value.

### HasDescriptor

`func (o *TaskModulr) HasDescriptor() bool`

HasDescriptor returns a boolean if a field has been set.

### GetStatus

`func (o *TaskModulr) GetStatus() PaymentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *TaskModulr) GetStatusOk() (*PaymentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *TaskModulr) SetStatus(v PaymentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *TaskModulr) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetState

`func (o *TaskModulr) GetState() map[string]interface{}`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *TaskModulr) GetStateOk() (*map[string]interface{}, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *TaskModulr) SetState(v map[string]interface{})`

SetState sets State field to given value.

### HasState

`func (o *TaskModulr) HasState() bool`

HasState returns a boolean if a field has been set.

### GetError

`func (o *TaskModulr) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *TaskModulr) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *TaskModulr) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *TaskModulr) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


