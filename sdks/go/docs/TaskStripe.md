# TaskStripe

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ConnectorID** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**Descriptor** | Pointer to [**TaskStripeDescriptor**](TaskStripeDescriptor.md) |  | [optional] 
**Status** | Pointer to [**PaymentStatus**](PaymentStatus.md) |  | [optional] 
**State** | Pointer to **map[string]interface{}** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewTaskStripe

`func NewTaskStripe() *TaskStripe`

NewTaskStripe instantiates a new TaskStripe object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaskStripeWithDefaults

`func NewTaskStripeWithDefaults() *TaskStripe`

NewTaskStripeWithDefaults instantiates a new TaskStripe object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TaskStripe) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaskStripe) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaskStripe) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TaskStripe) HasId() bool`

HasId returns a boolean if a field has been set.

### GetConnectorID

`func (o *TaskStripe) GetConnectorID() string`

GetConnectorID returns the ConnectorID field if non-nil, zero value otherwise.

### GetConnectorIDOk

`func (o *TaskStripe) GetConnectorIDOk() (*string, bool)`

GetConnectorIDOk returns a tuple with the ConnectorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorID

`func (o *TaskStripe) SetConnectorID(v string)`

SetConnectorID sets ConnectorID field to given value.

### HasConnectorID

`func (o *TaskStripe) HasConnectorID() bool`

HasConnectorID returns a boolean if a field has been set.

### GetCreatedAt

`func (o *TaskStripe) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *TaskStripe) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *TaskStripe) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *TaskStripe) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *TaskStripe) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *TaskStripe) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *TaskStripe) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *TaskStripe) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetDescriptor

`func (o *TaskStripe) GetDescriptor() TaskStripeDescriptor`

GetDescriptor returns the Descriptor field if non-nil, zero value otherwise.

### GetDescriptorOk

`func (o *TaskStripe) GetDescriptorOk() (*TaskStripeDescriptor, bool)`

GetDescriptorOk returns a tuple with the Descriptor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptor

`func (o *TaskStripe) SetDescriptor(v TaskStripeDescriptor)`

SetDescriptor sets Descriptor field to given value.

### HasDescriptor

`func (o *TaskStripe) HasDescriptor() bool`

HasDescriptor returns a boolean if a field has been set.

### GetStatus

`func (o *TaskStripe) GetStatus() PaymentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *TaskStripe) GetStatusOk() (*PaymentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *TaskStripe) SetStatus(v PaymentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *TaskStripe) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetState

`func (o *TaskStripe) GetState() map[string]interface{}`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *TaskStripe) GetStateOk() (*map[string]interface{}, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *TaskStripe) SetState(v map[string]interface{})`

SetState sets State field to given value.

### HasState

`func (o *TaskStripe) HasState() bool`

HasState returns a boolean if a field has been set.

### GetError

`func (o *TaskStripe) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *TaskStripe) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *TaskStripe) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *TaskStripe) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


